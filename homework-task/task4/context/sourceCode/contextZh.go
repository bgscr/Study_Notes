// Copyright 2014 The Go Authors. All rights reserved.
// 版权所有 2014 The Go Authors。保留所有权利。
// Use of this source code is governed by a BSD-style
// 本源码的使用受 BSD 风格许可证保护，许可证内容见 LICENSE 文件。
// license that can be found in the LICENSE file.

// Package context defines the Context type, which carries deadlines,
// 本包定义了 Context 类型，用于在 API 边界和进程之间传递截止时间、
// cancellation signals, and other request-scoped values across API boundaries
// 取消信号以及其他请求作用域内的值。
// and between processes.
//
// Incoming requests to a server should create a [Context], and outgoing
// 服务器收到的每个请求都应该创建一个 [Context]，而对外发起调用时则应接受一个 Context。
// calls to servers should accept a Context. The chain of function
// 它们之间的函数调用链必须将该 Context 向下传递，并可选择使用 [WithCancel]、
// calls between them must propagate the Context, optionally replacing
// [WithDeadline]、[WithTimeout] 或 [WithValue] 派生出的 Context 进行替换。
// it with a derived Context created using [WithCancel], [WithDeadline],
// [WithTimeout], or [WithValue].
//
// A Context may be canceled to indicate that work done on its behalf should stop.
// Context 可以被取消，以表明应停止代表它执行的工作。
// A Context with a deadline is canceled after the deadline passes.
// 带截止时间的 Context 在截止时间到达后会被取消。
// When a Context is canceled, all Contexts derived from it are also canceled.
// 一旦某个 Context 被取消，所有从它派生出的 Context 也会被取消。
//
// The [WithCancel], [WithDeadline], and [WithTimeout] functions take a
// [WithCancel]、[WithDeadline] 和 [WithTimeout] 函数接收一个
// Context (the parent) and return a derived Context (the child) and a [CancelFunc].
// Context（父）并返回一个派生的 Context（子）和一个 [CancelFunc]。
// Calling the CancelFunc directly cancels the child and its children,
// 直接调用 CancelFunc 会取消子 Context 及其子代，并移除父 Context 对该子的引用，
// removes the parent's reference to the child, and stops any associated timers.
// 同时停止所有关联的计时器。如果未调用 CancelFunc，则会造成子 Context 及其子代的泄露，
// Failing to call the CancelFunc leaks the child and its children until the parent is canceled.
// 直到父 Context 被取消为止。go vet 工具会检查所有控制流路径上是否正确使用了 CancelFunc。
// The go vet tool checks that CancelFuncs are used on all control-flow paths.
//
// The [WithCancelCause], [WithDeadlineCause], and [WithTimeoutCause] functions
// [WithCancelCause]、[WithDeadlineCause] 和 [WithTimeoutCause] 函数
// return a [CancelCauseFunc], which takes an error and records it as the cancellation cause.
// 返回一个 [CancelCauseFunc]，它接收一个 error 并将其记录为取消原因。
// Calling [Cause] on the canceled context or any of its children retrieves the cause.
// 在已取消的 Context 或其任意子 Context 上调用 [Cause] 可检索该原因。
// If no cause is specified, Cause(ctx) returns the same value as ctx.Err().
// 如果未指定原因，Cause(ctx) 返回与 ctx.Err() 相同的值。
//
// Programs that use Contexts should follow these rules to keep interfaces
// 使用 Context 的程序应遵循以下规则，以保持包间接口的一致性，
// consistent across packages and enable static analysis tools to check context
// 并使静态分析工具能够检查 Context 的传播：
// propagation:
//
// Do not store Contexts inside a struct type; instead, pass a Context
// 不要将 Context 存储在结构体类型内部；而是将 Context 显式地
// explicitly to each function that needs it. This is discussed further in
// 传递给每个需要它的函数。这一点在
// <https://go.dev/blog/context-and-structs> 中有进一步讨论。
// <https://go.dev/blog/context-and-structs>. The Context should be the first
// Context 应为第一个参数，通常命名为 ctx：
// parameter, typically named ctx:
//
//	func DoSomething(ctx context.Context, arg Arg) error {
//		// ... use ctx ...
//	}
//
// Do not pass a nil [Context], even if a function permits it. Pass [context.TODO]
// 即使函数允许，也不要传递 nil [Context]。如果不确定使用哪个 Context，
// if you are unsure about which Context to use.
// 请使用 [context.TODO]。
//
// Use context Values only for request-scoped data that transits processes and
// Context 的值仅应用于跨进程和 API 边界传递的请求作用域数据，
// APIs, not for passing optional parameters to functions.
// 而不应用于向函数传递可选参数。
//
// The same Context may be passed to functions running in different goroutines;
// 同一个 Context 可以传递给在不同 goroutine 中运行的函数；
// Contexts are safe for simultaneous use by multiple goroutines.
// Context 可被多个 goroutine 并发安全地使用。
//
// See <https://go.dev/blog/context> for example code for a server that uses
// 有关使用 Context 的服务器示例代码，请参见 <https://go.dev/blog/context>。
// Contexts.
package task4_context

import (
	"errors"
	"internal/reflectlite"
	"sync"
	"sync/atomic"
	"time"
)

// A Context carries a deadline, a cancellation signal, and other values across
// Context 携带截止时间、取消信号以及其他跨 API 边界的值。
// API boundaries.
//
// Context's methods may be called by multiple goroutines simultaneously.
// Context 的方法可以被多个 goroutine 并发调用。
type Context interface {
	// Deadline returns the time when work done on behalf of this context
	// Deadline 返回代表此 Context 执行的工作应被取消的时间。
	// should be canceled. Deadline returns ok==false when no deadline is
	// 当未设置截止时间时，Deadline 返回 ok==false。
	// set. Successive calls to Deadline return the same results.
	// 连续调用 Deadline 将返回相同的结果。
	Deadline() (deadline time.Time, ok bool)

	// Done returns a channel that's closed when work done on behalf of this
	// Done 返回一个 channel，当代表此 Context 的工作应被取消时该 channel 会被关闭。
	// context should be canceled. Done may return nil if this context can
	// 如果此 Context 永远不会被取消，Done 可能返回 nil。
	// never be canceled. Successive calls to Done return the same value.
	// 连续调用 Done 会返回相同的值。
	// The close of the Done channel may happen asynchronously,
	// Done channel 的关闭可能在 cancel 函数返回之后异步发生。
	// after the cancel function returns.
	//
	// WithCancel arranges for Done to be closed when cancel is called;
	// WithCancel 会在调用 cancel 时安排关闭 Done；
	// WithDeadline arranges for Done to be closed when the deadline
	// WithDeadline 会在截止时间到达时安排关闭 Done；
	// expires; WithTimeout arranges for Done to be closed when the timeout
	// WithTimeout 会在超时到达时安排关闭 Done。
	// elapses.
	//
	// Done is provided for use in select statements:
	// Done 用于在 select 语句中使用：
	//
	//  // Stream generates values with DoSomething and sends them to out
	//  // Stream 通过 DoSomething 生成值并发送到 out，
	//  // until DoSomething returns an error or ctx.Done is closed.
	//  // 直到 DoSomething 返回错误或 ctx.Done 被关闭。
	//  func Stream(ctx context.Context, out chan<- Value) error {
	//  	for {
	//  		v, err := DoSomething(ctx)
	//  		if err != nil {
	//  			return err
	//  		}
	//  		select {
	//  		case <-ctx.Done():
	//  			return ctx.Err()
	//  		case out <- v:
	//  		}
	//  	}
	//  }
	//
	// See https://blog.golang.org/pipelines for more examples of how to use
	// 更多关于如何使用 Done channel 进行取消的示例，请参见 https://blog.golang.org/pipelines。
	// a Done channel for cancellation.
	Done() <-chan struct{}

	// If Done is not yet closed, Err returns nil.
	// 如果 Done 尚未关闭，Err 返回 nil。
	// If Done is closed, Err returns a non-nil error explaining why:
	// 如果 Done 已关闭，Err 返回一个非 nil 的 error 说明原因：
	// DeadlineExceeded if the context's deadline passed,
	// DeadlineExceeded 表示 Context 的截止时间已过，
	// or Canceled if the context was canceled for some other reason.
	// 或者 Canceled 表示 Context 因其他原因被取消。
	// After Err returns a non-nil error, successive calls to Err return the same error.
	// 一旦 Err 返回非 nil error，后续调用 Err 将返回相同的 error。
	Err() error

	// Value returns the value associated with this context for key, or nil
	// Value 返回与 key 关联的值，若无关联值则返回 nil。
	// if no value is associated with key. Successive calls to Value with
	// 连续使用相同 key 调用 Value 会返回相同结果。
	// the same key returns the same result.
	//
	// Use context values only for request-scoped data that transits
	// Context 的值仅应用于跨进程和 API 边界传递的请求作用域数据，
	// processes and API boundaries, not for passing optional parameters to
	// 而非用于向函数传递可选参数。
	// functions.
	//
	// A key identifies a specific value in a Context. Functions that wish
	// key 用于在 Context 中标识特定值。希望在 Context 中存储值的函数
	// to store values in Context typically allocate a key in a global
	// 通常在全局变量中分配一个 key，然后将其用于 context.WithValue 和
	// variable then use that key as the argument to context.WithValue and
	// Context.Value 的参数。key 可以是任何支持相等比较的类型；
	// Context.Value. A key can be any type that supports equality;
	// 包应使用未导出类型作为 key 的定义，以避免冲突。
	// packages should define keys as an unexported type to avoid
	// collisions.
	//
	// Packages that define a Context key should provide type-safe accessors
	// 定义 Context key 的包应提供类型安全的访问器
	// for the values stored using that key:
	// 用于访问通过该 key 存储的值：
	//
	// 	// Package user defines a User type that's stored in Contexts.
	// 	// 包 user 定义了存储在 Context 中的 User 类型。
	// 	package user
	//
	// 	import "context"
	//
	// 	// User is the type of value stored in the Contexts.
	// 	// User 是存储在 Context 中的值的类型。
	// 	type User struct {...}
	//
	// 	// key is an unexported type for keys defined in this package.
	// 	// key 是本包定义的未导出 key 类型。
	// 	// This prevents collisions with keys defined in other packages.
	// 	// 这防止与其他包定义的 key 冲突。
	// 	type key int
	//
	// 	// userKey is the key for user.User values in Contexts. It is
	// 	// userKey 是 Context 中 user.User 值的 key。它是未导出的；
	// 	// unexported; clients use user.NewContext and user.FromContext
	// 	// 客户端应使用 user.NewContext 和 user.FromContext，而非直接使用此 key。
	// 	// instead of using this key directly.
	// 	var userKey key
	//
	// 	// NewContext returns a new Context that carries value u.
	// 	// NewContext 返回一个新的携带值 u 的 Context。
	// 	func NewContext(ctx context.Context, u *User) context.Context {
	// 		return context.WithValue(ctx, userKey, u)
	// 	}
	//
	// 	// FromContext returns the User value stored in ctx, if any.
	// 	// FromContext 返回 ctx 中存储的 User 值（若有）。
	// 	func FromContext(ctx context.Context) (*User, bool) {
	// 		u, ok := ctx.Value(userKey).(*User)
	// 		return u, ok
	// 	}
	Value(key any) any
}

// Canceled is the error returned by [Context.Err] when the context is canceled
// Canceled 是当 Context 因非截止时间原因被取消时，[Context.Err] 返回的 error。
// for some reason other than its deadline passing.
var Canceled = errors.New("context canceled")

// DeadlineExceeded is the error returned by [Context.Err] when the context is canceled
// DeadlineExceeded 是当 Context 因截止时间到达而被取消时，[Context.Err] 返回的 error。
// due to its deadline passing.
var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }

// An emptyCtx is never canceled, has no values, and has no deadline.
// emptyCtx 永远不会被取消，不携带值，也没有截止时间。
// It is the common base of backgroundCtx and todoCtx.
// 它是 backgroundCtx 和 todoCtx 的公共基础。
type emptyCtx struct{}

func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (emptyCtx) Done() <-chan struct{} {
	return nil
}

func (emptyCtx) Err() error {
	return nil
}

func (emptyCtx) Value(key any) any {
	return nil
}

type backgroundCtx struct{ emptyCtx }

func (backgroundCtx) String() string {
	return "context.Background"
}

type todoCtx struct{ emptyCtx }

func (todoCtx) String() string {
	return "context.TODO"
}

// Background returns a non-nil, empty [Context]. It is never canceled, has no
// Background 返回一个非 nil 的空的 [Context]。它永远不会被取消，
// values, and has no deadline. It is typically used by the main function,
// 不携带值，也没有截止时间。通常由 main 函数、初始化和测试使用，
// initialization, and tests, and as the top-level Context for incoming
// 并作为传入请求的顶层 Context。
// requests.
func Background() Context {
	return backgroundCtx{}
}

// TODO returns a non-nil, empty [Context]. Code should use context.TODO when
// TODO 返回一个非 nil 的空的 [Context]。当不清楚应使用哪个 Context，
// it's unclear which Context to use or it is not yet available (because the
// 或 Context 尚不可用（因为周围函数尚未扩展为接受 Context 参数）时，应使用 context.TODO。
// surrounding function has not yet been extended to accept a Context
// parameter).
func TODO() Context {
	return todoCtx{}
}

// A CancelFunc tells an operation to abandon its work.
// CancelFunc 通知操作放弃其工作。
// A CancelFunc does not wait for the work to stop.
// CancelFunc 不会等待工作停止。
// A CancelFunc may be called by multiple goroutines simultaneously.
// CancelFunc 可被多个 goroutine 同时调用。
// After the first call, subsequent calls to a CancelFunc do nothing.
// 第一次调用后，再次调用 CancelFunc 不会生效。
type CancelFunc func()

// WithCancel returns a derived context that points to the parent context
// WithCancel 返回一个派生 context，指向父 context，
// but has a new Done channel. The returned context's Done channel is closed
// 但拥有一个新的 Done channel。当返回的 cancel 函数被调用，
// when the returned cancel function is called or when the parent context's
// 或父 context 的 Done channel 被关闭时（以先发生者为准），
// Done channel is closed, whichever happens first.
// 返回的 context 的 Done channel 将被关闭。
//
// Canceling this context releases resources associated with it, so code should
// 取消此 context 会释放其关联资源，因此代码应在运行于此 [Context] 中的操作完成后尽快调用 cancel。
// call cancel as soon as the operations running in this [Context] complete.
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := withCancel(parent)
	return c, func() { c.cancel(true, Canceled, nil) }
}

// A CancelCauseFunc behaves like a [CancelFunc] but additionally sets the cancellation cause.
// CancelCauseFunc 的行为类似 [CancelFunc]，但还会设置取消原因。
// This cause can be retrieved by calling [Cause] on the canceled Context or on
// 该原因可通过在已取消的 Context 或其任意派生 Context 上调用 [Cause] 获取。
// any of its derived Contexts.
//
// If the context has already been canceled, CancelCauseFunc does not set the cause.
// 如果 Context 已被取消，CancelCauseFunc 不会设置原因。
// For example, if childContext is derived from parentContext:
// 例如，若 childContext 派生自 parentContext：
//   - if parentContext is canceled with cause1 before childContext is canceled with cause2,
//     如果 parentContext 在 childContext 以 cause2 取消之前以 cause1 取消，
//     then Cause(parentContext) == Cause(childContext) == cause1
//     则 Cause(parentContext) == Cause(childContext) == cause1
//   - if childContext is canceled with cause2 before parentContext is canceled with cause1,
//     如果 childContext 在 parentContext 以 cause1 取消之前以 cause2 取消，
//     then Cause(parentContext) == cause1 and Cause(childContext) == cause2
//     则 Cause(parentContext) == cause1 且 Cause(childContext) == cause2
type CancelCauseFunc func(cause error)

// WithCancelCause behaves like [WithCancel] but returns a [CancelCauseFunc] instead of a [CancelFunc].
// WithCancelCause 的行为类似 [WithCancel]，但返回的是 [CancelCauseFunc] 而非 [CancelFunc]。
// Calling cancel with a non-nil error (the "cause") records that error in ctx;
// 使用非 nil error（“cause”）调用 cancel 会将该 error 记录到 ctx 中；
// it can then be retrieved using Cause(ctx).
// 之后可通过 Cause(ctx) 获取。
// Calling cancel with nil sets the cause to Canceled.
// 使用 nil 调用 cancel 会将原因设为 Canceled。
//
// Example use:
// 示例用法：
//
//	ctx, cancel := context.WithCancelCause(parent)
//	cancel(myError)
//	ctx.Err() // returns context.Canceled
//	context.Cause(ctx) // returns myError
func WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc) {
	c := withCancel(parent)
	return c, func(cause error) { c.cancel(true, Canceled, cause) }
}

func withCancel(parent Context) *cancelCtx {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := &cancelCtx{}
	c.propagateCancel(parent, c)
	return c
}

// Cause returns a non-nil error explaining why c was canceled.
// Cause 返回一个非 nil error，解释 c 为何被取消。
// The first cancellation of c or one of its parents sets the cause.
// c 或其任意父级的首次取消设置原因。
// If that cancellation happened via a call to CancelCauseFunc(err),
// 如果该取消是通过调用 CancelCauseFunc(err) 发生的，
// then [Cause] returns err.
// 则 [Cause] 返回 err。
// Otherwise Cause(c) returns the same value as c.Err().
// 否则 Cause(c) 返回与 c.Err() 相同的值。
// Cause returns nil if c has not been canceled yet.
// 若 c 尚未被取消，Cause 返回 nil。
func Cause(c Context) error {
	if cc, ok := c.Value(&cancelCtxKey).(*cancelCtx); ok {
		cc.mu.Lock()
		cause := cc.cause
		cc.mu.Unlock()
		if cause != nil {
			return cause
		}
		// Either this context is not canceled,
		// 要么此 Context 未被取消，
		// or it is canceled and the cancellation happened in a
		// 要么已被取消，但取消发生在自定义 Context 实现而非 *cancelCtx 中。
		// custom context implementation rather than a *cancelCtx.
	}
	// There is no cancelCtxKey value with a cause, so we know that c is
	// 不存在带原因的 cancelCtxKey 值，因此可知 c 并非由 WithCancelCause 创建的已取消 Context 的后代。
	// not a descendant of some canceled Context created by WithCancelCause.
	// Therefore, there is no specific cause to return.
	// 因此，没有特定原因可返回。
	// If this is not one of the standard Context types,
	// 若这不是标准 Context 类型之一，
	// it might still have an error even though it won't have a cause.
	// 它仍可能有 error，尽管不会有 cause。
	return c.Err()
}

// AfterFunc arranges to call f in its own goroutine after ctx is canceled.
// AfterFunc 安排在 ctx 被取消后，在其自己的 goroutine 中调用 f。
// If ctx is already canceled, AfterFunc calls f immediately in its own goroutine.
// 若 ctx 已被取消，AfterFunc 会立即在其 goroutine 中调用 f。
//
// Multiple calls to AfterFunc on a context operate independently;
// 对同一 context 的多次 AfterFunc 调用相互独立；
// one does not replace another.
// 后者不会替代前者。
//
// Calling the returned stop function stops the association of ctx with f.
// 调用返回的 stop 函数会终止 ctx 与 f 的关联。
// It returns true if the call stopped f from being run.
// 若调用阻止了 f 的执行，则返回 true。
// If stop returns false,
// 若 stop 返回 false，
// either the context is canceled and f has been started in its own goroutine;
// 则要么 context 已被取消且 f 已在其 goroutine 中启动；
// or f was already stopped.
// 要么 f 已被停止。
// The stop function does not wait for f to complete before returning.
// stop 函数在返回前不会等待 f 完成。
// If the caller needs to know whether f is completed,
// 若调用者需要知道 f 是否完成，
// it must coordinate with f explicitly.
// 必须显式与 f 协调。
//
// If ctx has a "AfterFunc(func()) func() bool" method,
// 若 ctx 拥有 "AfterFunc(func()) func() bool" 方法，
// AfterFunc will use it to schedule the call.
// AfterFunc 将使用它来调度调用。
func AfterFunc(ctx Context, f func()) (stop func() bool) {
	a := &afterFuncCtx{
		f: f,
	}
	a.cancelCtx.propagateCancel(ctx, a)
	return func() bool {
		stopped := false
		a.once.Do(func() {
			stopped = true
		})
		if stopped {
			a.cancel(true, Canceled, nil)
		}
		return stopped
	}
}

type afterFuncer interface {
	AfterFunc(func()) func() bool
}

type afterFuncCtx struct {
	cancelCtx
	once sync.Once // either starts running f or stops f from running
	f    func()
}

func (a *afterFuncCtx) cancel(removeFromParent bool, err, cause error) {
	a.cancelCtx.cancel(false, err, cause)
	if removeFromParent {
		removeChild(a.Context, a)
	}
	a.once.Do(func() {
		go a.f()
	})
}

// A stopCtx is used as the parent context of a cancelCtx when
// 当 AfterFunc 已在父 context 注册时，stopCtx 用作 cancelCtx 的父 context。
// an AfterFunc has been registered with the parent.
// It holds the stop function used to unregister the AfterFunc.
// 它保存用于注销 AfterFunc 的 stop 函数。
type stopCtx struct {
	Context
	stop func() bool
}

// goroutines counts the number of goroutines ever created; for testing.
// goroutines 统计曾创建的 goroutine 数量，仅用于测试。
var goroutines atomic.Int32

// &cancelCtxKey is the key that a cancelCtx returns itself for.
// &cancelCtxKey 是 cancelCtx 返回自身的 key。
var cancelCtxKey int

// parentCancelCtx returns the underlying *cancelCtx for parent.
// parentCancelCtx 返回 parent 对应的底层 *cancelCtx。
// It does this by looking up parent.Value(&cancelCtxKey) to find
// 它通过查找 parent.Value(&cancelCtxKey) 来寻找
// the innermost enclosing *cancelCtx and then checking whether
// 最内层包裹的 *cancelCtx，然后检查
// parent.Done() matches that *cancelCtx. (If not, the *cancelCtx
// parent.Done() 是否与该 *cancelCtx 匹配。（若不匹配，则该 *cancelCtx
// has been wrapped in a custom implementation providing a
// 已被包裹在自定义实现中，提供了不同的 done channel，此时不应绕过它。）
// different done channel, in which case we should not bypass it.)
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	pdone, _ := p.done.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return p, true
}

// removeChild removes a context from its parent.
// removeChild 从父 context 中移除一个子 context。
func removeChild(parent Context, child canceler) {
	if s, ok := parent.(stopCtx); ok {
		s.stop()
		return
	}
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}

// A canceler is a context type that can be canceled directly. The
// canceler 是一个可直接取消的 context 类型。
// implementations are *cancelCtx and *timerCtx.
// 其实现为 *cancelCtx 和 *timerCtx。
type canceler interface {
	cancel(removeFromParent bool, err, cause error)
	Done() <-chan struct{}
}

// closedchan is a reusable closed channel.
// closedchan 是一个可复用的已关闭 channel。
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

// A cancelCtx can be canceled. When canceled, it also cancels any children
// cancelCtx 可被取消。取消时，它也会取消所有实现了 canceler 的子项。
// that implement canceler.
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      atomic.Value          // set to non-nil by the first cancel call
	cause    error                 // set to non-nil by the first cancel call
}

func (c *cancelCtx) Value(key any) any {
	if key == &cancelCtxKey {
		return c
	}
	return value(c.Context, key)
}

func (c *cancelCtx) Done() <-chan struct{} {
	d := c.done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.done.Load()
	if d == nil {
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}

func (c *cancelCtx) Err() error {
	// An atomic load is ~5x faster than a mutex, which can matter in tight loops.
	// 原子加载比互斥锁快约 5 倍，在紧凑循环中可能很重要。
	if err := c.err.Load(); err != nil {
		return err.(error)
	}
	return nil
}

// propagateCancel arranges for child to be canceled when parent is.
// propagateCancel 安排在 parent 被取消时取消 child。
// It sets the parent context of cancelCtx.
// 它设置 cancelCtx 的父 context。
func (c *cancelCtx) propagateCancel(parent Context, child canceler) {
	c.Context = parent

	done := parent.Done()
	if done == nil {
		return // parent is never canceled
	}

	select {
	case <-done:
		// parent is already canceled
		// parent 已被取消
		child.cancel(false, parent.Err(), Cause(parent))
		return
	default:
	}

	if p, ok := parentCancelCtx(parent); ok {
		// parent is a *cancelCtx, or derives from one.
		// parent 是 *cancelCtx 或派生自 *cancelCtx。
		p.mu.Lock()
		if err := p.err.Load(); err != nil {
			// parent has already been canceled
			// parent 已被取消
			child.cancel(false, err.(error), p.cause)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
		return
	}

	if a, ok := parent.(afterFuncer); ok {
		// parent implements an AfterFunc method.
		// parent 实现了 AfterFunc 方法。
		c.mu.Lock()
		stop := a.AfterFunc(func() {
			child.cancel(false, parent.Err(), Cause(parent))
		})
		c.Context = stopCtx{
			Context: parent,
			stop:    stop,
		}
		c.mu.Unlock()
		return
	}

	goroutines.Add(1)
	go func() {
		select {
		case <-parent.Done():
			child.cancel(false, parent.Err(), Cause(parent))
		case <-child.Done():
		}
	}()
}

type stringer interface {
	String() string
}

func contextName(c Context) string {
	if s, ok := c.(stringer); ok {
		return s.String()
	}
	return reflectlite.TypeOf(c).String()
}

func (c *cancelCtx) String() string {
	return contextName(c.Context) + ".WithCancel"
}

// cancel closes c.done, cancels each of c's children, and, if
// cancel 关闭 c.done，取消 c 的每个子项，并且若 removeFromParent 为 true，
// removeFromParent is true, removes c from its parent's children.
// 则将 c 从父的 children 中移除。
// cancel sets c.cause to cause if this is the first time c is canceled.
// 若这是首次取消 c，cancel 将 c.cause 设置为 cause。
func (c *cancelCtx) cancel(removeFromParent bool, err, cause error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	if cause == nil {
		cause = err
	}
	c.mu.Lock()
	if c.err.Load() != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err.Store(err)
	c.cause = cause
	d, _ := c.done.Load().(chan struct{})
	if d == nil {
		c.done.Store(closedchan)
	} else {
		close(d)
	}
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		// 注意：在持有父级锁的同时获取子级的锁。
		child.cancel(false, err, cause)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}

// WithoutCancel returns a derived context that points to the parent context
// WithoutCancel 返回一个派生 context，指向父 context，
// and is not canceled when parent is canceled.
// 且在父被取消时不会被取消。
// The returned context returns no Deadline or Err, and its Done channel is nil.
// 返回的 context 不返回 Deadline 或 Err，其 Done channel 为 nil。
// Calling [Cause] on the returned context returns nil.
// 在返回的 context 上调用 [Cause] 返回 nil。
func WithoutCancel(parent Context) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	return withoutCancelCtx{parent}
}

type withoutCancelCtx struct {
	c Context
}

func (withoutCancelCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (withoutCancelCtx) Done() <-chan struct{} {
	return nil
}

func (withoutCancelCtx) Err() error {
	return nil
}

func (c withoutCancelCtx) Value(key any) any {
	return value(c, key)
}

func (c withoutCancelCtx) String() string {
	return contextName(c.c) + ".WithoutCancel"
}

// WithDeadline returns a derived context that points to the parent context
// WithDeadline 返回一个派生 context，指向父 context，
// but has the deadline adjusted to be no later than d. If the parent's
// 但截止时间调整为不晚于 d。若父的截止时间已早于 d，
// deadline is already earlier than d, WithDeadline(parent, d) is semantically
// 则 WithDeadline(parent, d) 在语义上等价于 parent。
// equivalent to parent. The returned [Context.Done] channel is closed when
// 返回的 [Context.Done] channel 会在截止时间到达、
// the deadline expires, when the returned cancel function is called,
// 返回的 cancel 函数被调用，或父 context 的 Done channel 被关闭时关闭，
// or when the parent context's Done channel is closed, whichever happens first.
// 以先发生者为准。
//
// Canceling this context releases resources associated with it, so code should
// 取消此 context 会释放其关联资源，因此代码应在运行于此 [Context] 中的操作完成后尽快调用 cancel。
// call cancel as soon as the operations running in this [Context] complete.
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	return WithDeadlineCause(parent, d, nil)
}

// WithDeadlineCause behaves like [WithDeadline] but also sets the cause of the
// WithDeadlineCause 的行为类似 [WithDeadline]，但还会在截止超时时设置返回 Context 的原因。
// returned Context when the deadline is exceeded. The returned [CancelFunc] does
// 返回的 [CancelFunc] 不会设置原因。
// not set the cause.
func WithDeadlineCause(parent Context, d time.Time, cause error) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		// The current deadline is already sooner than the new one.
		// 当前截止时间已早于新截止时间。
		return WithCancel(parent)
	}
	c := &timerCtx{
		deadline: d,
	}
	c.cancelCtx.propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded, cause) // deadline has already passed
		return c, func() { c.cancel(false, Canceled, nil) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err.Load() == nil {
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded, cause)
		})
	}
	return c, func() { c.cancel(true, Canceled, nil) }
}

// A timerCtx carries a timer and a deadline. It embeds a cancelCtx to
// timerCtx 携带计时器和截止时间。它内嵌 cancelCtx 以实现 Done 和 Err。
// implement Done and Err. It implements cancel by stopping its timer then
// 它通过停止计时器然后委托给 cancelCtx.cancel 来实现取消。
// delegating to cancelCtx.cancel.
type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) String() string {
	return contextName(c.cancelCtx.Context) + ".WithDeadline(" +
		c.deadline.String() + " [" +
		time.Until(c.deadline).String() + "])"
}

func (c *timerCtx) cancel(removeFromParent bool, err, cause error) {
	c.cancelCtx.cancel(false, err, cause)
	if removeFromParent {
		// Remove this timerCtx from its parent cancelCtx's children.
		// 从父 cancelCtx 的 children 中移除此 timerCtx。
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}

// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
//
// Canceling this context releases resources associated with it, so code should
// 取消此 context 会释放其关联资源，因此代码应在运行于此 [Context] 中的操作完成后尽快调用 cancel。
// call cancel as soon as the operations running in this [Context] complete:
//
//	func slowOperationWithTimeout(ctx context.Context) (Result, error) {
//		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
//		defer cancel()  // releases resources if slowOperation completes before timeout elapses
//		// 若 slowOperation 在超时前完成，则释放资源
//		return slowOperation(ctx)
//	}
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

// WithTimeoutCause behaves like [WithTimeout] but also sets the cause of the
// WithTimeoutCause 的行为类似 [WithTimeout]，但还会在超时时设置返回 Context 的原因。
// returned Context when the timeout expires. The returned [CancelFunc] does
// 返回的 [CancelFunc] 不会设置原因。
// not set the cause.
func WithTimeoutCause(parent Context, timeout time.Duration, cause error) (Context, CancelFunc) {
	return WithDeadlineCause(parent, time.Now().Add(timeout), cause)
}

// WithValue returns a derived context that points to the parent Context.
// WithValue 返回一个派生 context，指向父 Context。
// In the derived context, the value associated with key is val.
// 在派生 context 中，key 关联的值为 val。
//
// Use context Values only for request-scoped data that transits processes and
// Context 的值仅应用于跨进程和 API 边界传递的请求作用域数据，
// APIs, not for passing optional parameters to functions.
// 而非用于向函数传递可选参数。
//
// The provided key must be comparable and should not be of type
// 提供的 key 必须可比较，且不应为 string 或任何其他内置类型，
// string or any other built-in type to avoid collisions between
// 以避免使用 context 的包之间发生冲突。
// packages using context. Users of WithValue should define their own
// WithValue 的使用者应定义自己的 key 类型。
// types for keys. To avoid allocating when assigning to an
// 为避免赋值给 interface{} 时分配内存，
// interface{}, context keys often have concrete type
// context key 通常为具体类型 struct{}。
// struct{}. Alternatively, exported context key variables' static
// 或者，导出的 context key 变量的静态类型应为指针或 interface。
// type should be a pointer or interface.
func WithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}

// A valueCtx carries a key-value pair. It implements Value for that key and
// valueCtx 携带一个键值对。它为该 key 实现 Value，并将其他所有调用委托给内嵌的 Context。
// delegates all other calls to the embedded Context.
type valueCtx struct {
	Context
	key, val any
}

// stringify tries a bit to stringify v, without using fmt, since we don't
// stringify 尝试将 v 转为字符串，不使用 fmt，因为不想依赖 unicode 表。
// want context depending on the unicode tables. This is only used by
// 此函数仅用于 *valueCtx.String()。
// *valueCtx.String().
func stringify(v any) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case string:
		return s
	case nil:
		return "<nil>"
	}
	return reflectlite.TypeOf(v).String()
}

func (c *valueCtx) String() string {
	return contextName(c.Context) + ".WithValue(" +
		stringify(c.key) + ", " +
		stringify(c.val) + ")"
}

func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	return value(c.Context, key)
}

func value(c Context, key any) any {
	for {
		switch ctx := c.(type) {
		case *valueCtx:
			if key == ctx.key {
				return ctx.val
			}
			c = ctx.Context
		case *cancelCtx:
			if key == &cancelCtxKey {
				return c
			}
			c = ctx.Context
		case withoutCancelCtx:
			if key == &cancelCtxKey {
				// This implements Cause(ctx) == nil
				// 这使得 Cause(ctx) == nil
				// when ctx is created using WithoutCancel.
				// 当 ctx 使用 WithoutCancel 创建时。
				return nil
			}
			c = ctx.c
		case *timerCtx:
			if key == &cancelCtxKey {
				return &ctx.cancelCtx
			}
			c = ctx.Context
		case backgroundCtx, todoCtx:
			return nil
		default:
			return c.Value(key)
		}
	}
}

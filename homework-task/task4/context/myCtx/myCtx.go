package myCtx

import (
	"context"
	"sync"
	"time"
)

type MyCancelCtx struct {
	parent   context.Context
	mu       sync.Mutex
	done     chan struct{}
	err      error
	children map[context.Context]func() // 保存子节点 cancel 函数
}

func newMyCancelCtx(parent context.Context) *MyCancelCtx {
	return &MyCancelCtx{
		parent:   parent,
		done:     make(chan struct{}),
		children: make(map[context.Context]func()),
	}
}

func (c *MyCancelCtx) Deadline() (time.Time, bool) { return c.parent.Deadline() }
func (c *MyCancelCtx) Done() <-chan struct{}       { return c.done }
func (c *MyCancelCtx) Err() error                  { return c.err }
func (c *MyCancelCtx) Value(key any) any           { return c.parent.Value(key) }

func (c *MyCancelCtx) cancel(err error) {
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return
	}
	c.err = err
	close(c.done)
	for _, childCancel := range c.children {
		childCancel()
	}
	c.children = nil
	c.mu.Unlock()
}

func WithMyCancel(parent context.Context) (context.Context, context.CancelFunc) {
	child := newMyCancelCtx(parent)
	go func() {
		select {
		case <-parent.Done():
			child.cancel(parent.Err())
		case <-child.Done():
		}
	}()
	return child, func() { child.cancel(context.Canceled) }
}

type valueCtx struct {
	context.Context
	key, val any
}

func (c *valueCtx) Value(key any) any {
	if key == c.key {
		return c.val
	}
	return c.Context.Value(key)
}

func WithMyValue(parent context.Context, key, val any) context.Context {
	return &valueCtx{parent, key, val}
}

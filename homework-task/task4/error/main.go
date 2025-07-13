package main

import "fmt"

type myError struct {
	code    int // 小写字母私有化
	message string
}

func (e *myError) Error() string {
	return fmt.Sprintf("错误码:%d, 消息:%s", e.code, e.message)
}

func NewMyError(code int, msg string) error {
	return &myError{code: code, message: msg}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			if err, ok := r.(*myError); ok {
				fmt.Println("recover 捕获异常:", err)
				return
			} else {
				panic(r)
			}
		}
	}()

	panic(NewMyError(500, "自定义错误"))

}

// Pakcage coroutine 协程逻辑封装
package coroutine

import (
	"fmt"
	"runtime"
	"sync"

	"google.golang.org/grpc/grpclog"
)

// ErrPanic 协程运行中发生Panic
var ErrPanic = fmt.Errorf("panic found in call handlers")

// panicBufLen panic调用栈日志buffer大小，默认1024
var panicBufLen = 1024

// Recover 程序panic后处理
func Recover() {
	if r := recover(); r != nil {
		buf := make([]byte, panicBufLen)
		buf = buf[:runtime.Stack(buf, false)]
		grpclog.Errorf("[PANIC]%v\n%s\n", r, buf)
	}
}

// Go 创建协程自带恢复机制
func Go(f func()) {
	go func() {
		defer Recover()
		f()
	}()
}

// GoAndWait 并发启动多个协程，并等待所有协程返回
func GoAndWait(handlers ...func() error) error {
	var (
		wg   sync.WaitGroup
		once sync.Once
		err  error
	)
	for _, f := range handlers {
		wg.Add(1)
		go func(handler func() error) {
			defer func() {
				if e := recover(); e != nil {
					buf := make([]byte, panicBufLen)
					buf = buf[:runtime.Stack(buf, false)]
					grpclog.Errorf("[PANIC]%v\n%s\n", e, buf)
					once.Do(func() {
						err = ErrPanic
					})
				}
				wg.Done()
			}()
			if e := handler(); e != nil {
				once.Do(func() {
					err = e
				})
			}
		}(f)
	}
	wg.Wait()
	return err
}

// GoAndWaitWithConcurrency 批量rpc调用，并发数concurrency
func GoAndWaitWithConcurrency(concurrency int, handles []func() error) error {
	var (
		wg      sync.WaitGroup
		once    sync.Once // 保护返回值err，只需要赋值一次即可
		err     error
		conChan = make(chan struct{}, concurrency) // 控制并发量
	)
	for _, h := range handles {
		conChan <- struct{}{}
		wg.Add(1)
		go func(handle func() error) {
			defer func() {
				if e := recover(); e != nil {
					buf := make([]byte, 1024*10)
					buf = buf[:runtime.Stack(buf, false)]
					grpclog.Errorf("[PANIC]%v\n%s\n", e, buf)
					once.Do(func() {
						err = ErrPanic
					})
				}
				<-conChan
				wg.Done()
			}()
			if e := handle(); e != nil {
				once.Do(func() {
					err = e
				})
			}
		}(h)
	}
	wg.Wait()
	return err
}

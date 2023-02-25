// handle 函数逻辑封装
package handle

import (
	"time"
)

// DoWithRetry 运行handle函数，发生错误时重试
func DoWithRetry(h func() error, retryCount int, interval time.Duration) error {
	alwayRetry := func(e error) bool { return true }
	return DoWithRetryWhenSomeErrs(h, alwayRetry, retryCount, interval)
}

// DoWithRetryWhenSomeErrs 运行handle函数，仅特定错误时重试
func DoWithRetryWhenSomeErrs(h func() error, needRetry func(error) bool, retryCount int, interval time.Duration) error {
	if interval < 0 {
		interval = time.Millisecond * 200 // 默认重试时间间隔
	}
	var e error
	for i := 1; i <= retryCount; i++ {
		if e = h(); e == nil {
			return nil
		}
		if i == retryCount || !needRetry(e) { // 超过重试次数或者无需重试
			return e
		}
		time.Sleep(interval)
	}
	return e
}

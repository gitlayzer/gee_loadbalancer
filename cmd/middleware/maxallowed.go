package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

// MaxAllowedMiddleware 用于限制并发请求数量
func MaxAllowedMiddleware(n uint) mux.MiddlewareFunc {
	sem := make(chan struct{}, n)           // 创建一个信号量，容量为 n
	acquire := func() { sem <- struct{}{} } // 申请一个信号量
	release := func() { <-sem }             // 释放一个信号量

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acquire()            // 申请一个信号量
			defer release()      // 释放一个信号量
			next.ServeHTTP(w, r) // 执行下一个中间件
		})
	}
}

package handler

import (
	"log"
	"net/http"
)


//拦截器
func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
	
	return http.HandlerFunc(
		func(w http.ResponseWriter,r *http.Request) {
			log.Println("token拦截器.....")
			userName := r.FormValue("username")
			token := r.FormValue("token")

			if len(userName) < 3 || !IsTokenValid(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h(w,r)
		},
	)
}
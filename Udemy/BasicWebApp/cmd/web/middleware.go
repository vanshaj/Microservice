package main

import (
	"fmt"
	"net/http"
	"os"
)

func SetCookie(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		cookies := &http.Cookie{
			Name:     "t",
			Value:    "token",
			HttpOnly: true,
			Path:     "/",
			Secure:   false,
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(w, cookies)
		fmt.Fprintln(os.Stdout, "Hit the api")
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func SessionLoad(next http.Handler) http.Handler {
	return app.SessionManager.LoadAndSave(next)
}

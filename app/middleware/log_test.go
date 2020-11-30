package middleware

import (
	. "github.com/milkbobo/fishgoweb/app/log"
	. "github.com/milkbobo/fishgoweb/app/router"
	"net/http"
	"testing"
)

func TestLog(t *testing.T) {
	log, _ := NewLog(LogConfig{
		Driver: "console",
	})

	factory := NewRouterFactory()
	factory.Use(NewLogMiddleware(log, nil))
	factory.GET("/a", func(w http.ResponseWriter, r *http.Request) {

	})
	factory.GET("/b", func(w http.ResponseWriter, r *http.Request) {
		panic("Hello World")
	})

	router := factory.Create()
	testCase := []string{"/a", "/b", "/c"}
	for _, url := range testCase {
		r, _ := http.NewRequest("GET", url, nil)
		w := &fakeWriter{}
		router.ServeHTTP(w, r)
	}

}

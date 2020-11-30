package middleware

import (
	"context"
	"errors"
	. "github.com/milkbobo/fishgoweb/app/log"
	. "github.com/milkbobo/fishgoweb/app/render"
	. "github.com/milkbobo/fishgoweb/app/router"
	. "github.com/milkbobo/fishgoweb/app/session"
	. "github.com/milkbobo/fishgoweb/app/validator"
	. "github.com/milkbobo/fishgoweb/assert"
	. "github.com/milkbobo/fishgoweb/encoding"
	. "github.com/milkbobo/fishgoweb/language"
	"net/http"
	"testing"
)

func a_json(v Validator, s Session) interface{} {
	return v.MustQuery("key")
}

func b_Json(v Validator, s Session) interface{} {
	Throw(10001, "my god")
	return nil
}

func c(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Fish"))
}

func d_Json(v Validator, s Session) interface{} {
	return errors.New("my god2")
}

func e_Json(v Validator, s Session) interface{} {
	return NewException(10002, "my god%v", 3)
}

type dStruct struct {
}

func (this *dStruct) f_Json(v Validator, s Session) interface{} {
	return "my god4"
}

func (this *dStruct) G_Json(v Validator, s Session) interface{} {
	return "my god5"
}

func (this *dStruct) h_Json(ctx context.Context, v Validator, s Session) interface{} {
	return ctx.Value("db")
}

func jsonToArray(data string) interface{} {
	var result interface{}
	err := DecodeJson([]byte(data), &result)
	if err != nil {
		panic(err)
	}
	return result
}

func TestEasy(t *testing.T) {
	log, _ := NewLog(LogConfig{Driver: "console"})
	renderFactory, _ := NewRenderFactory(RenderConfig{})
	validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
	sessionFactory, _ := NewSessionFactory(SessionConfig{Driver: "memory", CookieName: "fishmm"})
	middleware := NewEasyMiddleware(log, validatorFactory, sessionFactory, renderFactory, nil)

	factory := NewRouterFactory()
	d := &dStruct{}
	factory.Use(NewLogMiddleware(log, nil))
	factory.Use(middleware)
	factory.GET("/a", a_json)
	factory.GET("/b", b_Json)
	factory.GET("/c", c)
	factory.GET("/d", d_Json)
	factory.GET("/e", e_Json)
	factory.GET("/f", d.f_Json)
	factory.GET("/g", d.G_Json)
	router := factory.Create()

	r, _ := http.NewRequest("GET", "http://example.com/a?key=mmc", nil)
	w := &fakeWriter{}
	router.ServeHTTP(w, r)

	AssertEqual(t, jsonToArray(w.Read()), map[string]interface{}{"code": 0.0, "msg": "", "data": "mmc"})

	r2, _ := http.NewRequest("GET", "http://example.com/b?key2=mmc", nil)
	w2 := &fakeWriter{}
	router.ServeHTTP(w2, r2)
	AssertEqual(t, jsonToArray(w2.Read()), map[string]interface{}{"code": 10001.0, "msg": "my god", "data": nil})

	r3, _ := http.NewRequest("GET", "http://example.com/c", nil)
	w3 := &fakeWriter{}
	router.ServeHTTP(w3, r3)
	AssertEqual(t, w3.Read(), "Hello Fish")

	r4, _ := http.NewRequest("GET", "http://example.com/d", nil)
	w4 := &fakeWriter{}
	router.ServeHTTP(w4, r4)
	AssertEqual(t, jsonToArray(w4.Read()), map[string]interface{}{"code": 1.0, "msg": "my god2", "data": nil})

	r5, _ := http.NewRequest("GET", "http://example.com/e", nil)
	w5 := &fakeWriter{}
	router.ServeHTTP(w5, r5)
	AssertEqual(t, jsonToArray(w5.Read()), map[string]interface{}{"code": 10002.0, "msg": "my god3", "data": nil})

	r6, _ := http.NewRequest("GET", "http://example.com/f", nil)
	w6 := &fakeWriter{}
	router.ServeHTTP(w6, r6)
	AssertEqual(t, jsonToArray(w6.Read()), map[string]interface{}{"code": 0.0, "data": "my god4", "msg": ""})

	r7, _ := http.NewRequest("GET", "http://example.com/g", nil)
	w7 := &fakeWriter{}
	router.ServeHTTP(w7, r7)
	AssertEqual(t, jsonToArray(w7.Read()), map[string]interface{}{"code": 0.0, "data": "my god5", "msg": ""})
}

func TestEasy2(t *testing.T) {
	log, _ := NewLog(LogConfig{Driver: "console"})
	renderFactory, _ := NewRenderFactory(RenderConfig{})
	validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
	sessionFactory, _ := NewSessionFactory(SessionConfig{Driver: "memory", CookieName: "fishmm"})
	middleware := NewEasyMiddleware(log, validatorFactory, sessionFactory, renderFactory, nil)

	factory := NewRouterFactory()
	factory.Use(NewLogMiddleware(log, nil))
	factory.Use(middleware)
	factory.Use(func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		lastHandler, isOk := prev.Handler.(func(ctx context.Context, v Validator, s Session) interface{})
		if isOk == false {
			return prev
		}
		return RouterMiddlewareContext{
			Data: prev.Data,
			Handler: func(v Validator, s Session) interface{} {
				newContext := context.Background()
				newContext = context.WithValue(newContext, "db", "contextDbValue")
				return lastHandler(newContext, v, s)
			},
		}
	})

	d := &dStruct{}
	factory.GET("/a", d.h_Json)
	router := factory.Create()

	r, _ := http.NewRequest("GET", "http://example.com/a", nil)
	w := &fakeWriter{}
	router.ServeHTTP(w, r)
	AssertEqual(t, jsonToArray(w.Read()), map[string]interface{}{"code": 0.0, "msg": "", "data": "contextDbValue"})
}

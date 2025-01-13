package gee

import(
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct{
	router map[string]HandlerFunc	//路由映射表
}

//New构造gee.Engine
func New() *Engine{
	return &Engine{router: make(map[string]HandlerFunc)}
}

//向路由映射表中添加路由及处理方法
func (engine *Engine) addRoute(method string, pattern string, handle HandlerFunc){
	key := method + "-" + pattern
	engine.router[key] = handle
}

//添加GET、POST请求的路由
func (engine *Engine) GET(pattern string, handle HandlerFunc){
	engine.addRoute("GET", pattern, handle)
}

func (engine *Engine) POST(pattern string, handle HandlerFunc){
	engine.addRoute("POST", pattern, handle)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr,engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok{
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
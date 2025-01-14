package gee

import(
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct{
	router 	*Router//路由
}

//New构造gee.Engine
func New() *Engine{
	return &Engine{router: newRouter()}
}

//添加GET、POST请求的路由
func (engine *Engine) GET(pattern string, handler HandlerFunc){
	engine.router.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc){
	engine.router.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr,engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	c := newContext(w,req)	//新上下文
	engine.router.handle(c)
}
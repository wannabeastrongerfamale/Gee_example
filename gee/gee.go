package gee

import(
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct{
	*RouterGroup	//根RouterGroup
	router 	*Router	//路由
	groups []*RouterGroup	//分组列表
}

type RouterGroup struct{
	prefix string
	engine *Engine
	middlewares []HandlerFunc
}

//New构造gee.Engine
func New() *Engine{
	engine := &Engine{router: newRouter()}
	//创建一个默认路由组--prefix为空
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//Group嵌套
func (group *RouterGroup) Group(prefix string) *RouterGroup{
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc){
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

//添加GET、POST请求的路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc){
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc){
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr,engine)
}

// 处理HTTP请求入口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	//在接收到一个具体请求时，判断请求适用于哪些中间件
	var middlewares []HandlerFunc
	for _, group := range engine.groups{
		if strings.HasPrefix(req.URL.Path, group.prefix){
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w,req)	//新上下文
	c.handlers = middlewares
	engine.router.handle(c)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc){
	group.middlewares = append(group.middlewares, middlewares...)
}

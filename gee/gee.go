package gee

import(
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct{
	*RouterGroup	//嵌入RouterGroup---可以直接使用RouterGroup的成员方法
	router 	*Router	//路由
	groups []*RouterGroup	//分组列表
}

type RouterGroup struct{
	prefix string
	parent *RouterGroup
	engine *Engine
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
		parent: group,
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

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	c := newContext(w,req)	//新上下文
	engine.router.handle(c)
}
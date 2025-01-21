package gee

import(
	"net/http"
	"strings"
	"html/template"
	"path"
)

type HandlerFunc func(c *Context)

type Engine struct{
	*RouterGroup	//根RouterGroup
	router 	*Router	//路由
	groups []*RouterGroup	//分组列表
	htmlTemplates *template.Template // for html render
	funcMap template.FuncMap
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

// 创建static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	//创建HTTP文件处理器，包装器StripPrefix从请求的URL中删除给定前缀，并在fs目录下寻找文件
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

//添加静态文件路由
func (group *RouterGroup) Static(relativePath string, root string){
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	//fmt.Printf("%q", urlPattern)

	group.GET(urlPattern, handler)
}

//自定义模板渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap){
	engine.funcMap = funcMap
}

//解析路径下的所有文件并与模板类关联
func (engine *Engine) LoadHTMLGlob(pattern string){
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

//添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc){
	group.middlewares = append(group.middlewares, middlewares...)
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

	// 实例化新上下文
	c := newContext(w, req, middlewares, engine)
	engine.router.handle(c)
}

//启动服务器
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr,engine)
}
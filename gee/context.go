package gee

import(
	"net/http"
	"fmt"
	"encoding/json"
)

type H map[string]interface{}	//别名

type Context struct{
	//	origin objects
	Writer http.ResponseWriter
	Req *http.Request
	// request info
	Path string
	Method string
	// 动态路由参数
	Params map[string]string
	// respone info
	StatusCode int
	// middlewares info
	handlers []HandlerFunc
	index int
	
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request, middlewares []HandlerFunc, engine *Engine) *Context{
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		handlers: middlewares,
		index: -1,
		engine: engine,
	}
}

// 获取URL上指定的请求参数值-GET请求
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 获取消息体中的指定键的键值-POST请求
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Status(code int){
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string){
	c.Writer.Header().Set(key,value)
}

//提供快速构造String/Date/JSON/HTML响应的方式
func (c *Context) String(code int, format string, values ...interface{}){
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}){
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}){
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	//执行特定模板，并写回c.Writer
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil{
		c.Fail(500, err.Error())
	}
}

func (c *Context) Fail(code int, data string){
	c.Status(code)
	c.Writer.Write([]byte(data))
}

func (c *Context) Next(){
	/*
	// 1、简单递归--仅适用与所有中间件中都有next()
	c.index++
	c.handlers[index](c)
	*/
	//2、高适配模式--适配中间件含或不含next()
	c.index++
	for ; c.index < len(c.handlers); c.index++{
		c.handlers[c.index](c)	//handle
	}
}
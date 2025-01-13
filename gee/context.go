package gee

import(
	"net/http"
	"fmt"
	"encoding/json"
)

type H map[string]interface{}

type Context struct{
	//	origin objects
	Writer http.ResponseWriter
	Req *http.Request
	// request info
	Path string
	Method string
	// respone info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context{
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
	}
}

// 获取URL上指定的请求参数值-GET请求
func (c *Context) Query(key string) string{
	return c.Req.URL.Query().Get(key)
}

// 获取消息体中的指定键的键值-POST请求
func (c *Context) PostForm(key string) string{
	return c.Req.FormValue(key)
}

func (c *Context) Status(code int){
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string){
	c.Writer.Header().Set(key,value)
}

func (c *Context) String(code int, format string, values ...interface{}){
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, json string){
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string){
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	c.Writer.Write([]byte(html))
}
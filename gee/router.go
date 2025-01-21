package gee

import(
	"fmt"
	"strings"
	"net/http"
)

type Router struct{
	handlers map[string]HandlerFunc	//路由映射表
	roots map[string]*node
}
// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *Router{
	return &Router{handlers: make(map[string]HandlerFunc), roots: make(map[string]*node)}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//注册路由
func (r *Router) addRoute(method string, pattern string, handler HandlerFunc){
	parts := parsePattern(pattern)

	//在前缀树中插入路由
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)

	//在路由方法映射表中注册
	key := method + "-" + pattern
	r.handlers[key] = handler
}

//匹配路由
func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	//在前缀树中匹配路由
	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

//根据路由映射表进行handle
func (r *Router) handle(c *Context){
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil{
		key := c.Method + "-" + node.pattern
		handler := r.handlers[key]
		fmt.Printf("%q\n", node.pattern)	//输出当前HTTP请求路由
		c.Params = params
		c.handlers = append(c.handlers, handler)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	
	// 开始handle下一个处理函数
	c.Next()
}
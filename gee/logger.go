package gee

import(
	"log"
	"time"
	"os"
)

func Logger(fileDirectory string, fileName string) HandlerFunc {
	return func(c *Context) {
		filePath := fileDirectory + fileName
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)	//log输出默认写入标准错误输出stderr
		}
		//将log的输出重定向到filePath
		log.SetOutput(file)
		
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
# limiter
Golang实现的HTTP客户端下载限速功能包

示例代码：
```Go
package main

import (
	"log"
	"net/http"

	"github.com/dxvgef/limiter"
)

func main() {
	http.HandleFunc("/", func (resp http.ResponseWriter, req *http.Request) {
		// 传输demo.mp4文件，限速每秒100KKB
         if err := limiter.ServeFile(resp, req, "./demo.mp4", 100*1024); err != nil {
         	resp.WriteHeader(500)
            resp.Write([]byte(err.Error()))
         }
     })

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err.Error())
        return
	}
}
```
# limiter
Golang实现的HTTP限流功能包

示例代码：
```Go
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dxvgef/limiter"
)

func serveFile(resp http.ResponseWriter, req *http.Request) {
    err := limiter.ServeFile(resp, req, "./demo.mp4", 100*1024)
    if err != nil {
    	resp.WriteHeader(500)
    	resp.Write([]byte(err.Error()))
    }
}

func main() {
	// log.SetFlags(log.Lshortfile)

	http.HandleFunc("/", serveFile)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
```
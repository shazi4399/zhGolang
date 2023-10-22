package main

/**
https = http + tls
首先，生成证书和key。
go run $GOROOT/src/crypto/tls/generate_cert.go --host="localhost"
go run "D:/Program Files/Go/src/crypto/tls/generate_cert.go" --host="localhost"
在当前目录下会生成cert.pem和key.pem这两个文件。我把这两个文件移到了config/keys目录下(这个操作不是必须的)
*/

import (
	"net/http"

	"github.com/unrolled/secure" //需要用到这个第三方库
)

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("欢迎来到HTTP Secure的世界"))
})

func main1() {
	secureMiddleware := secure.New(secure.Options{
		//把http://localhost:5678重定向到https://localhost:5678。这个选项其实可以不写，它是默认行为
		SSLRedirect: true,
		SSLHost:     "localhost:5678",
	})
	app := secureMiddleware.Handler(myHandler)
	//启动https（http+tls）服务
	http.ListenAndServeTLS("localhost:5678", "config/keys/cert.pem", "config/keys/key.pem", app)
}

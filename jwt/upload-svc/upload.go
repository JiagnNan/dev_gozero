// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net/http"

	"upload-svc/internal/config"
	"upload-svc/internal/handler"
	"upload-svc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/upload-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// WithFileServer 把 UploadDir 挂到 /files 前缀下，
	// 这样 logic 返回的 /files/xxx 才能真正被浏览器访问到。
	// 官方文档只给了 URL 拼接、没给这段注册，照抄会得到 404。
	server := rest.MustNewServer(c.RestConf,
		rest.WithFileServer("/files", http.Dir(c.UploadDir)),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

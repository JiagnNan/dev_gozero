// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"jwt/internal/config"
	"jwt/internal/handler"
	"jwt/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// TODO: 接入 GORM 后传入真实仓储，例如：
	//   db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//   ctx := svc.NewServiceContext(c, model.NewUserModel(db))
	// 目前传 nil，login / profile 等依赖 UserModel 的接口会 panic，
	// 但 logout（纯 Redis 操作）链路可独立运行与验证。
	ctx := svc.NewServiceContext(c, nil)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

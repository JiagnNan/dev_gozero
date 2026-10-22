// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	// UploadDir 文件落盘目录。
	// 用相对路径，以进程启动目录（upload-svc）为基准，换机器时不用改配置。
	UploadDir string
}

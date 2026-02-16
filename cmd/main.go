package main

import (
	"fmt"
	"os"
	"strconv"
	"test/configs"
	"test/cron"
	"test/dao"
	"test/pkg/email"
	"test/router"
)

func main() {
	if err := configs.Init(); err != nil {
		panic("配置加载失败: " + err.Error())
	}
	fmt.Println("✅ 配置加载成功")

	if err := dao.InitDB(); err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	fmt.Println("✅ 数据库连接成功")

	email.Init()
	cron.Init()
	defer cron.Stop()

	r := router.SetupRouter()

	port := configs.Cfg.Server.Port

	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	fmt.Printf("服务启动，监听端口: %d\n", port)
	fmt.Printf("API 地址: http://localhost:%d\n", port)

	r.Run(fmt.Sprintf(":%d", port))
}

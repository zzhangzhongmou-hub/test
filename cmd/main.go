package main

import (
	"fmt"
	"test/configs"
	"test/dao"
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

	r := router.SetupRouter()

	port := configs.Cfg.Server.Port
	fmt.Printf("服务启动，监听端口: %d\n", port)
	fmt.Printf("API 地址: http://localhost:%d\n", port)
	fmt.Println("测试登录: POST /user/login")
	fmt.Println("测试个人信息: GET /user/profile (需要 Authorization Header)")

	r.Run(fmt.Sprintf(":%d", port))
}

// main.go - 最简单的启动文件
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin实例
	r := gin.Default()

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "人事调动管理系统API服务正常运行",
		})
	})

	// 欢迎页面
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "欢迎使用人事调动管理系统API",
			"version": "1.0.0",
			"endpoints": []string{
				"GET /api/employees - 获取员工列表",
				"GET /api/employees/:id - 获取员工详情",
				"POST /api/employees - 创建员工",
				"POST /api/auth/login - 用户登录",
			},
		})
	})

	// 启动服务器
	port := ":8080"
	log.Printf("🚀 API服务启动在 http://localhost%s", port)
	log.Println("📋 可用接口：")
	log.Println("  GET  /health          - 健康检查")
	log.Println("  GET  /                - API文档")

	if err := r.Run(port); err != nil {
		log.Fatal("启动失败:", err)
	}
}

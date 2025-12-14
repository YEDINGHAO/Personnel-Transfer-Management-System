// main.go - 完整版
package main

import (
	"log"

	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/api"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	db, err := database.Init()
	if err != nil {
		log.Printf("⚠️  数据库连接失败，但API服务将继续运行（仅使用内存数据）: %v", err)
	} else {
		log.Println("✅ 数据库连接成功")
		// 测试数据库连接
		var result int
		db.Raw("SELECT 1").Scan(&result)
		log.Printf("✅ 数据库连接测试成功: %d", result)
	}

	// 创建Gin实例
	r := gin.Default()

	// 添加CORS中间件（允许跨域请求）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 创建控制器实例
	employeeController := api.EmployeeController{}

	// API路由组
	apiGroup := r.Group("/api")
	{
		// 员工管理
		apiGroup.GET("/employees", employeeController.GetEmployees)
		apiGroup.GET("/employees/:id", employeeController.GetEmployee)
		apiGroup.POST("/employees", employeeController.CreateEmployee)
		apiGroup.PUT("/employees/:id", employeeController.UpdateEmployee)
		apiGroup.DELETE("/employees/:id", employeeController.DeleteEmployee)
	}

	// 基础路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":     0,
			"message":  "人事调动管理系统 API",
			"version":  "1.0.0",
			"database": err == nil,
			"endpoints": []gin.H{
				{"method": "GET", "path": "/api/employees", "description": "获取员工列表"},
				{"method": "GET", "path": "/api/employees/:id", "description": "获取员工详情"},
				{"method": "POST", "path": "/api/employees", "description": "创建员工"},
				{"method": "PUT", "path": "/api/employees/:id", "description": "更新员工"},
				{"method": "DELETE", "path": "/api/employees/:id", "description": "删除员工"},
			},
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"database": err == nil,
		})
	})

	// 启动服务器
	port := ":8080"
	log.Printf("\n🚀 服务器启动在 http://localhost%s", port)
	log.Println("\n📋 可用接口：")
	log.Println("  GET    /                    - API文档")
	log.Println("  GET    /health              - 健康检查")
	log.Println("  GET    /api/employees       - 获取员工列表")
	log.Println("  GET    /api/employees/:id   - 获取员工详情")
	log.Println("  POST   /api/employees       - 创建员工")
	log.Println("  PUT    /api/employees/:id   - 更新员工")
	log.Println("  DELETE /api/employees/:id   - 删除员工")
	log.Println("\n💡 提示：")
	log.Println("  使用 curl 或 Postman 测试API")
	log.Println("  创建员工示例：")
	log.Println(`    curl -X POST http://localhost:8080/api/employees \
      -H "Content-Type: application/json" \
      -d '{
        "employee_id": "EMP001",
        "name": "张三",
        "status": 1,
        "arrival_date": "2024-01-15",
        "job_title": "软件工程师",
        "position": "高级工程师",
        "department": "技术部",
        "phone": "13800138001",
        "email": "zhangsan@company.com"
      }'`)

	if err := r.Run(port); err != nil {
		log.Fatal("启动失败:", err)
	}
}

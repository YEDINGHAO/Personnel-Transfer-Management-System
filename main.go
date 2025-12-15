// main.go
package main

import (
	"log"

	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/api"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/database"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化数据库
	db, err := database.Init()
	if err != nil {
		log.Printf("⚠️ 数据库连接失败: %v", err)
	} else {
		// 自动迁移所有模型表 (确保包含 Transfer 和 Department)
		db.AutoMigrate(&models.User{}, &models.Employee{}, &models.Department{}, &models.Transfer{})
		log.Println("✅ 数据库表结构同步完成")
	}

	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")                                   //允许所有域名（*）访问你的 API
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")    //允许的 HTTP 方法
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Token") //允许客户端携带的请求头
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 实例化控制器
	authCtrl := api.AuthController{}
	empCtrl := api.EmployeeController{}
	deptCtrl := api.DepartmentController{} // 新增
	transCtrl := api.TransferController{}  // 新增
	backupCtrl := api.BackupController{}   // 新增

	apiGroup := r.Group("/api")
	{
		// --- 认证模块 ---
		apiGroup.POST("/login", authCtrl.Login)
		apiGroup.POST("/register", authCtrl.Register) // 开发测试用

		// --- 员工管理模块 ---
		apiGroup.GET("/employees", empCtrl.GetEmployees)
		apiGroup.GET("/employees/:id", empCtrl.GetEmployee)
		apiGroup.POST("/employees", empCtrl.CreateEmployee)
		apiGroup.PUT("/employees/:id", empCtrl.UpdateEmployee)
		apiGroup.DELETE("/employees/:id", empCtrl.DeleteEmployee)

		// --- 部门管理模块 (新增) ---
		apiGroup.GET("/departments", deptCtrl.GetDepartments)
		apiGroup.POST("/departments", deptCtrl.CreateDepartment)

		// --- 调动管理子系统 (新增) ---
		// 1. 提交调动/退休申请
		apiGroup.POST("/transfers", transCtrl.CreateTransfer)
		// 2. 获取调动记录列表 (可筛选待审批)
		apiGroup.GET("/transfers", transCtrl.GetTransfers)
		// 3. 审批调动 (通过后自动更新员工表)
		apiGroup.PUT("/transfers/:id/approve", transCtrl.ApproveTransfer)

		// --- 系统维护模块 (新增) ---
		// 导出员工数据备份
		apiGroup.GET("/backup/export", backupCtrl.ExportEmployees)
	}

	// 启动服务
	port := ":8080"
	log.Printf("🚀 服务器启动在 http://localhost%s", port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}

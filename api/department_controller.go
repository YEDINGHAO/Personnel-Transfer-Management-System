// api/department_controller.go
package api

import (
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/database"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/models"
	"github.com/gin-gonic/gin"
)

type DepartmentController struct{}

// GetDepartments 获取所有部门
func (dc *DepartmentController) GetDepartments(c *gin.Context) {
	db := database.GetDB()
	var depts []models.Department
	db.Find(&depts)
	success(c, depts)
}

// CreateDepartment 创建部门
func (dc *DepartmentController) CreateDepartment(c *gin.Context) {
	var req models.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "参数错误")
		return
	}

	dept := models.Department{
		DeptNo:    req.DeptNo,
		Name:      req.Name,
		ManagerID: req.ManagerID,
	}

	db := database.GetDB()
	if err := db.Create(&dept).Error; err != nil {
		errorResponse(c, 500, "创建部门失败")
		return
	}
	success(c, dept)
}

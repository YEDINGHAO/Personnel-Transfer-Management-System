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

// UpdateDepartment 更新部门
func (dc *DepartmentController) UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "参数错误")
		return
	}

	db := database.GetDB()
	var dept models.Department
	if err := db.First(&dept, id).Error; err != nil {
		errorResponse(c, 404, "部门不存在")
		return
	}

	dept.DeptNo = req.DeptNo
	dept.Name = req.Name
	dept.ManagerID = req.ManagerID

	if err := db.Save(&dept).Error; err != nil {
		errorResponse(c, 500, "更新部门失败")
		return
	}

	success(c, dept)
}

// DeleteDepartment 删除部门
func (dc *DepartmentController) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	if err := db.Delete(&models.Department{}, id).Error; err != nil {
		errorResponse(c, 500, "删除部门失败")
		return
	}
	success(c, gin.H{"message": "删除成功"})
}

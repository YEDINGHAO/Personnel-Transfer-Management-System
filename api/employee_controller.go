// api/employee_controller.go
package api

import (
	"net/http"
	"strconv"

	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/database"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/models"
	"github.com/gin-gonic/gin"
)

type EmployeeController struct{}

// 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 分页响应
type PaginatedResponse struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// 成功响应
func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// 错误响应
func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// GetEmployees 获取员工列表
// @Summary 获取员工列表
// @Description 获取员工列表，支持分页和筛选
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "员工姓名"
// @Param status query int false "员工状态"
// @Param department query string false "部门名称"
// @Success 200 {object} Response{data=PaginatedResponse}
// @Router /api/employees [get]
func (ec *EmployeeController) GetEmployees(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	name := c.Query("name")
	status := c.Query("status")
	department := c.Query("department")

	// 确保分页参数有效
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取数据库连接
	db := database.GetDB()

	// 构建查询
	query := db.Model(&models.Employee{})

	// 添加筛选条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if status != "" {
		if s, err := strconv.Atoi(status); err == nil && s >= 1 && s <= 6 {
			query = query.Where("status = ?", s)
		}
	}

	if department != "" {
		query = query.Where("department LIKE ?", "%"+department+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 查询数据
	var employees []models.Employee
	err := query.Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&employees).Error

	if err != nil {
		errorResponse(c, 500, "查询员工列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	employeeResponses := make([]models.EmployeeResponse, len(employees))
	for i, emp := range employees {
		employeeResponses[i] = models.EmployeeResponse{
			ID:          emp.ID,
			EmployeeID:  emp.EmployeeID,
			Name:        emp.Name,
			Status:      emp.Status,
			StatusText:  models.GetStatusText(emp.Status),
			ArrivalDate: emp.ArrivalDate,
			JobTitle:    emp.JobTitle,
			Position:    emp.Position,
			Department:  emp.Department,
			Phone:       emp.Phone,
			Email:       emp.Email,
			Address:     emp.Address,
			Remark:      emp.Remark,
			CreatedAt:   emp.CreatedAt,
			UpdatedAt:   emp.UpdatedAt,
		}
	}

	// 返回分页响应
	success(c, PaginatedResponse{
		Items:    employeeResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetEmployee 获取单个员工
// @Summary 获取员工详情
// @Description 根据ID获取员工详情
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Success 200 {object} Response{data=models.EmployeeResponse}
// @Router /api/employees/{id} [get]
func (ec *EmployeeController) GetEmployee(c *gin.Context) {
	id := c.Param("id")

	// 参数验证
	employeeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		errorResponse(c, 400, "无效的员工ID")
		return
	}

	db := database.GetDB()
	var employee models.Employee

	// 查询员工
	err = db.First(&employee, employeeID).Error
	if err != nil {
		errorResponse(c, 404, "员工不存在")
		return
	}

	// 转换为响应格式
	employeeResponse := models.EmployeeResponse{
		ID:          employee.ID,
		EmployeeID:  employee.EmployeeID,
		Name:        employee.Name,
		Status:      employee.Status,
		StatusText:  models.GetStatusText(employee.Status),
		ArrivalDate: employee.ArrivalDate,
		JobTitle:    employee.JobTitle,
		Position:    employee.Position,
		Department:  employee.Department,
		Phone:       employee.Phone,
		Email:       employee.Email,
		Address:     employee.Address,
		Remark:      employee.Remark,
		CreatedAt:   employee.CreatedAt,
		UpdatedAt:   employee.UpdatedAt,
	}

	success(c, employeeResponse)
}

// CreateEmployee 创建员工
// @Summary 创建员工
// @Description 创建新员工
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param request body models.CreateEmployeeRequest true "员工信息"
// @Success 200 {object} Response{data=models.EmployeeResponse}
// @Router /api/employees [post]
func (ec *EmployeeController) CreateEmployee(c *gin.Context) {
	var req models.CreateEmployeeRequest

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 验证员工编号是否已存在
	db := database.GetDB()
	var count int64
	db.Model(&models.Employee{}).Where("employee_id = ?", req.EmployeeID).Count(&count)
	if count > 0 {
		errorResponse(c, 400, "员工编号已存在")
		return
	}

	// 创建员工对象
	employee := models.Employee{
		EmployeeID:  req.EmployeeID,
		Name:        req.Name,
		Status:      req.Status,
		ArrivalDate: req.ArrivalDate,
		JobTitle:    req.JobTitle,
		Position:    req.Position,
		Department:  req.Department,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		Remark:      req.Remark,
	}

	// 保存到数据库
	result := db.Create(&employee)
	if result.Error != nil {
		errorResponse(c, 500, "创建员工失败: "+result.Error.Error())
		return
	}

	// 查询刚创建的员工（为了获取完整的字段）
	db.First(&employee, employee.ID)

	// 返回响应
	employeeResponse := models.EmployeeResponse{
		ID:          employee.ID,
		EmployeeID:  employee.EmployeeID,
		Name:        employee.Name,
		Status:      employee.Status,
		StatusText:  models.GetStatusText(employee.Status),
		ArrivalDate: employee.ArrivalDate,
		JobTitle:    employee.JobTitle,
		Position:    employee.Position,
		Department:  employee.Department,
		Phone:       employee.Phone,
		Email:       employee.Email,
		Address:     employee.Address,
		Remark:      employee.Remark,
		CreatedAt:   employee.CreatedAt,
		UpdatedAt:   employee.UpdatedAt,
	}

	success(c, employeeResponse)
}

// UpdateEmployee 更新员工
// @Summary 更新员工信息
// @Description 更新员工信息
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Param request body models.UpdateEmployeeRequest true "员工信息"
// @Success 200 {object} Response{data=models.EmployeeResponse}
// @Router /api/employees/{id} [put]
func (ec *EmployeeController) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateEmployeeRequest

	// 参数验证
	employeeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		errorResponse(c, 400, "无效的员工ID")
		return
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	var employee models.Employee

	// 查找员工
	err = db.First(&employee, employeeID).Error
	if err != nil {
		errorResponse(c, 404, "员工不存在")
		return
	}

	// 更新字段
	updateData := make(map[string]interface{})

	if req.Name != "" {
		updateData["name"] = req.Name
	}

	if req.Status > 0 && req.Status <= 6 {
		updateData["status"] = req.Status
	}

	if req.JobTitle != "" {
		updateData["job_title"] = req.JobTitle
	}

	if req.Position != "" {
		updateData["position"] = req.Position
	}

	if req.Department != "" {
		updateData["department"] = req.Department
	}

	if req.Phone != "" {
		updateData["phone"] = req.Phone
	}

	if req.Email != "" {
		updateData["email"] = req.Email
	}

	if req.Address != "" {
		updateData["address"] = req.Address
	}

	if req.Remark != "" {
		updateData["remark"] = req.Remark
	}

	// 更新数据库
	result := db.Model(&employee).Updates(updateData)
	if result.Error != nil {
		errorResponse(c, 500, "更新员工失败: "+result.Error.Error())
		return
	}

	// 重新查询更新后的员工
	db.First(&employee, employeeID)

	// 返回响应
	employeeResponse := models.EmployeeResponse{
		ID:          employee.ID,
		EmployeeID:  employee.EmployeeID,
		Name:        employee.Name,
		Status:      employee.Status,
		StatusText:  models.GetStatusText(employee.Status),
		ArrivalDate: employee.ArrivalDate,
		JobTitle:    employee.JobTitle,
		Position:    employee.Position,
		Department:  employee.Department,
		Phone:       employee.Phone,
		Email:       employee.Email,
		Address:     employee.Address,
		Remark:      employee.Remark,
		CreatedAt:   employee.CreatedAt,
		UpdatedAt:   employee.UpdatedAt,
	}

	success(c, employeeResponse)
}

// DeleteEmployee 删除员工
// @Summary 删除员工
// @Description 删除员工
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param id path int true "员工ID"
// @Success 200 {object} Response
// @Router /api/employees/{id} [delete]
func (ec *EmployeeController) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	// 参数验证
	employeeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		errorResponse(c, 400, "无效的员工ID")
		return
	}

	db := database.GetDB()

	// 检查员工是否存在
	var count int64
	db.Model(&models.Employee{}).Where("id = ?", employeeID).Count(&count)
	if count == 0 {
		errorResponse(c, 404, "员工不存在")
		return
	}

	// 删除员工
	result := db.Delete(&models.Employee{}, employeeID)
	if result.Error != nil {
		errorResponse(c, 500, "删除员工失败: "+result.Error.Error())
		return
	}

	success(c, gin.H{
		"message": "员工删除成功",
		"id":      employeeID,
	})
}

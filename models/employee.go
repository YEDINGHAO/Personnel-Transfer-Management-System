// models/employee.go
package models

import (
	"time"
)

// EmployeeStatus 员工状态枚举
type EmployeeStatus int

const (
	StatusActive    EmployeeStatus = 1 // 在职
	StatusPartTime  EmployeeStatus = 2 // 兼职
	StatusProbation EmployeeStatus = 3 // 试用
	StatusResigned  EmployeeStatus = 4 // 离职
	StatusRehired   EmployeeStatus = 5 // 返聘
	StatusRetired   EmployeeStatus = 6 // 退休
)

// Employee 员工模型
type Employee struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	EmployeeID  string    `gorm:"size:20;uniqueIndex;not null" json:"employee_id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Status      int       `gorm:"not null;default:1" json:"status"` // 使用int而不是EmployeeStatus，便于JSON序列化
	ArrivalDate string    `gorm:"type:date;not null" json:"arrival_date"`
	JobTitle    string    `gorm:"size:100" json:"job_title"`
	Position    string    `gorm:"size:100" json:"position"`
	Department  string    `gorm:"size:100" json:"department"` // 先简单处理，后面再关联
	Phone       string    `gorm:"size:20" json:"phone"`
	Email       string    `gorm:"size:100" json:"email"`
	Address     string    `gorm:"type:text" json:"address"`
	Remark      string    `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Employee) TableName() string {
	return "employees"
}

// 请求和响应DTO
type CreateEmployeeRequest struct {
	EmployeeID  string `json:"employee_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Status      int    `json:"status" binding:"required,min=1,max=6"`
	ArrivalDate string `json:"arrival_date" binding:"required"`
	JobTitle    string `json:"job_title"`
	Position    string `json:"position"`
	Department  string `json:"department"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Remark      string `json:"remark"`
}

type UpdateEmployeeRequest struct {
	Name       string `json:"name"`
	Status     int    `json:"status" binding:"min=1,max=6"`
	JobTitle   string `json:"job_title"`
	Position   string `json:"position"`
	Department string `json:"department"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	Remark     string `json:"remark"`
}

type EmployeeResponse struct {
	ID          uint      `json:"id"`
	EmployeeID  string    `json:"employee_id"`
	Name        string    `json:"name"`
	Status      int       `json:"status"`
	StatusText  string    `json:"status_text"`
	ArrivalDate string    `json:"arrival_date"`
	JobTitle    string    `json:"job_title"`
	Position    string    `json:"position"`
	Department  string    `json:"department"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 获取状态文本
func GetStatusText(status int) string {
	statusMap := map[int]string{
		1: "在职",
		2: "兼职",
		3: "试用",
		4: "离职",
		5: "返聘",
		6: "退休",
	}
	if text, ok := statusMap[status]; ok {
		return text
	}
	return "未知"
}

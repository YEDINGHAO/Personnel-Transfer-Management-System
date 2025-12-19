// models/department.go
package models

import "time"

type Department struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DeptNo    string    `gorm:"size:20;uniqueIndex;not null" json:"dept_no"`   // 部门编号
	Name      string    `gorm:"size:100;not null" json:"name"`                 // 部门名称
	ManagerID uint      `json:"manager_id"`                                    // 部门主管ID (关联员工)
	Manager   Employee  `gorm:"foreignKey:ManagerID" json:"manager,omitempty"` // 主管信息
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDepartmentRequest struct {
	DeptNo    string `json:"dept_no" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ManagerID uint   `json:"manager_id"`
}

type UpdateDepartmentRequest struct {
	DeptNo    string `json:"dept_no" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ManagerID uint   `json:"manager_id"`
}

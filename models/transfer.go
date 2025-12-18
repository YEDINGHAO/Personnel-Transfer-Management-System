// models/transfer.go
package models

import "time"

// TransferType 调动类型枚举
const (
	TransferTypeDepartment = 1 // 部门调动
	TransferTypePosition   = 2 // 职位调动
	TransferTypeRetirement = 3 // 离退休
)

// TransferStatus 调动状态枚举
const (
	TransferStatusPending   = 1 // 待审批
	TransferStatusApproved  = 2 // 已批准
	TransferStatusRejected  = 3 // 已驳回
	TransferStatusCompleted = 4 // 已完成
)

type Transfer struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	EmployeeID   uint       `gorm:"not null;index" json:"employee_id"`
	Employee     Employee   `gorm:"foreignKey:EmployeeID;references:ID;constraint:-" json:"employee"`
	TransferDate string     `gorm:"type:date;not null" json:"transfer_date"`
	Type         int        `gorm:"not null" json:"type"`
	FromDeptID   *uint      `json:"from_dept_id"`
	FromDept     Department `gorm:"foreignKey:FromDeptID" json:"from_dept"`
	ToDeptID     *uint      `json:"to_dept_id"`
	ToDept       Department `gorm:"foreignKey:ToDeptID" json:"to_dept"`
	Reason       string     `gorm:"type:text" json:"reason"`
	Status       int        `gorm:"not null;default:1" json:"status"`
	ApproverID   uint       `json:"approver_id"`
	Approver     User       `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

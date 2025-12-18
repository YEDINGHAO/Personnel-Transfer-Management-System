// api/transfer_controller.go
package api

import (
	"fmt"
	"time"

	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/database"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransferController struct{}

// CreateTransferRequest 创建调动申请请求
type CreateTransferRequest struct {
	EmployeeID   uint   `json:"employee_id" binding:"required"`
	Type         int    `json:"type" binding:"required"` // 1-部门调动, 2-职位调动, 3-离退休
	TransferDate string `json:"transfer_date" binding:"required"`
	FromDeptID   uint   `json:"from_dept_id"`
	ToDeptID     uint   `json:"to_dept_id"` // 如果是部门调动，必填
	Reason       string `json:"reason"`
}

// ApproveTransferRequest 审批请求
type ApproveTransferRequest struct {
	Status     int  `json:"status" binding:"required,oneof=2 3"` // 2-通过, 3-驳回
	ApproverID uint `json:"approver_id"`                         // 实际应从Token获取
}

// CreateTransfer 创建调动/离退休申请
func (tc *TransferController) CreateTransfer(c *gin.Context) {
	var req CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 验证员工是否存在
	var emp models.Employee
	if err := db.First(&emp, req.EmployeeID).Error; err != nil {
		errorResponse(c, 404, "员工不存在")
		return
	}

	// 处理部门ID，0 视为未设置
	var fromDeptID *uint
	if req.FromDeptID != 0 {
		fromDeptID = &req.FromDeptID
		// 校验调出部门是否存在
		var fromDept models.Department
		if err := db.First(&fromDept, req.FromDeptID).Error; err != nil {
			errorResponse(c, 400, "调出部门不存在")
			return
		}
	}

	var toDeptID *uint
	if req.ToDeptID != 0 {
		toDeptID = &req.ToDeptID
		// 校验调入部门是否存在
		var toDept models.Department
		if err := db.First(&toDept, req.ToDeptID).Error; err != nil {
			errorResponse(c, 400, "调入部门不存在")
			return
		}
	}

	transfer := models.Transfer{
		EmployeeID:   req.EmployeeID,
		Type:         req.Type,
		TransferDate: req.TransferDate,
		FromDeptID:   fromDeptID,
		ToDeptID:     toDeptID,
		Reason:       req.Reason,
		Status:       models.TransferStatusPending,
		CreatedAt:    time.Now(),
	}

	if err := db.Omit("ApproverID", "ApprovedAt").Create(&transfer).Error; err != nil {
		errorResponse(c, 500, "创建申请失败")
		return
	}

	success(c, transfer)
}

// GetTransfers 获取调动记录列表
func (tc *TransferController) GetTransfers(c *gin.Context) {
	employeeID := c.Query("employee_id")
	status := c.Query("status") // 1-待审批

	db := database.GetDB()
	query := db.Model(&models.Transfer{}).Preload("Employee").Preload("FromDept").Preload("ToDept")

	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var transfers []models.Transfer
	query.Order("created_at desc").Find(&transfers)

	success(c, transfers)
}

// ApproveTransfer 审批调动 (核心逻辑)
// 审批通过后，会自动更新 Employee 表的数据
func (tc *TransferController) ApproveTransfer(c *gin.Context) {
	id := c.Param("id")
	var req ApproveTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "参数错误")
		return
	}

	db := database.GetDB()
	var transfer models.Transfer
	if err := db.First(&transfer, id).Error; err != nil {
		errorResponse(c, 404, "调动记录不存在")
		return
	}

	if transfer.Status != models.TransferStatusPending {
		errorResponse(c, 400, "该记录已审批，无法重复操作")
		return
	}

	// 开启事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新调动表状态
		transfer.Status = req.Status
		transfer.ApproverID = req.ApproverID
		now := time.Now()
		transfer.ApprovedAt = &now

		if err := tx.Save(&transfer).Error; err != nil {
			return err
		}

		// 如果是驳回(3)，则不修改员工表，直接返回
		if req.Status == models.TransferStatusRejected {
			return nil
		}

		// 2. 如果是批准(2)，则修改员工基本表 (Requirement 2 & 3)
		var employee models.Employee
		if err := tx.First(&employee, transfer.EmployeeID).Error; err != nil {
			return err
		}

		// 根据调动类型更新员工信息
		if transfer.Type == models.TransferTypeDepartment {
			// 部门调动：查询新部门名称并更新
			if transfer.ToDeptID == nil {
				return fmt.Errorf("目标部门未设置")
			}
			var newDept models.Department
			if err := tx.First(&newDept, *transfer.ToDeptID).Error; err != nil {
				return fmt.Errorf("目标部门不存在")
			}
			employee.Department = newDept.Name // 更新部门
		} else if transfer.Type == models.TransferTypeRetirement {
			// 离退休：更新状态为退休(6)
			employee.Status = 6
		}

		// 保存员工变更
		if err := tx.Save(&employee).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		errorResponse(c, 500, "审批处理失败: "+err.Error())
		return
	}

	success(c, gin.H{"message": "审批完成，员工信息已同步更新"})
}

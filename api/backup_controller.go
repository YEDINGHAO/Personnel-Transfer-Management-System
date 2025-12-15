// api/backup_controller.go
package api

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/database"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/models"
	"github.com/gin-gonic/gin"
)

type BackupController struct{}

// ExportEmployees 导出员工信息为CSV
func (bc *BackupController) ExportEmployees(c *gin.Context) {
	db := database.GetDB()
	var employees []models.Employee
	if err := db.Find(&employees).Error; err != nil {
		errorResponse(c, 500, "获取数据失败")
		return
	}

	// 设置响应头，告诉浏览器这是一个下载文件
	fileName := fmt.Sprintf("employees_backup_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+fileName)

	// 创建 CSV Writer
	writer := csv.NewWriter(c.Writer)

	// 写入 UTF-8 BOM 防止 Excel 打开乱码
	c.Writer.Write([]byte("\xEF\xBB\xBF"))

	// 写入表头
	header := []string{"ID", "员工编号", "姓名", "状态", "部门", "职位", "入职日期", "电话", "邮箱"}
	writer.Write(header)

	// 写入数据行
	for _, emp := range employees {
		row := []string{
			fmt.Sprintf("%d", emp.ID),
			emp.EmployeeID,
			emp.Name,
			models.GetStatusText(emp.Status),
			emp.Department,
			emp.Position,
			emp.ArrivalDate,
			emp.Phone,
			emp.Email,
		}
		writer.Write(row)
	}

	writer.Flush()
}

// api/auth_controller.go
package api

import (
	"time"

	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/database"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/models"
	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

// Register 用户注册（开发测试用）
func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	db := database.GetDB()
	var count int64
	db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		errorResponse(c, 400, "用户名已存在")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errorResponse(c, 500, "密码加密失败")
		return
	}

	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		RealName:  req.RealName,
		Email:     req.Email,
		Phone:     req.Phone,
		LastLogin: time.Now(),
	}

	result := db.Create(&user)
	if result.Error != nil {
		errorResponse(c, 500, "注册失败")
		return
	}

	// 不返回密码
	user.Password = ""

	success(c, user)
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	db := database.GetDB()
	var user models.User

	// 查找用户
	err := db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		errorResponse(c, 401, "用户名或密码错误")
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		errorResponse(c, 401, "用户名或密码错误")
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		errorResponse(c, 500, "生成令牌失败")
		return
	}

	// 更新最后登录时间
	db.Model(&user).Update("last_login", time.Now())

	// 不返回密码
	user.Password = ""

	success(c, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile 获取当前用户信息
func (ac *AuthController) GetProfile(c *gin.Context) {
	// 从上下文获取用户ID（需要中间件）
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, 401, "未认证")
		return
	}

	db := database.GetDB()
	var user models.User

	err := db.First(&user, userID).Error
	if err != nil {
		errorResponse(c, 404, "用户不存在")
		return
	}

	// 不返回密码
	user.Password = ""

	success(c, user)
}

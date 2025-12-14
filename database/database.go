// database/database.go
package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/YEDINGHAO/Personnel-Transfer-Management-System/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Config 数据库配置
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Charset  string `json:"charset"`
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "123456", // 修改为你的密码
		Name:     "hrms",
		Charset:  "utf8mb4",
	}
}

// Init 初始化数据库连接
func Init() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		cfg := GetDefaultConfig()

		// 构建DSN
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.Charset)

		// 连接数据库
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			// 关闭默认事务，提高性能
			SkipDefaultTransaction: true,
		})

		if err != nil {
			err = fmt.Errorf("连接数据库失败: %v", err)
			return
		}

		log.Println("✅ 数据库连接成功")

		// 自动迁移表结构
		err = AutoMigrate()
		if err != nil {
			log.Printf("警告: 自动迁移失败: %v", err)
			// 不返回错误，继续运行（表可能已存在）
		}
	})

	return db, err
}

// AutoMigrate 自动迁移表结构
func AutoMigrate() error {
	// 导入模型
	models := []interface{}{
		&models.Employee{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移表失败: %v", err)
		}
	}

	log.Println("✅ 数据库表结构迁移完成")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		_, err := Init()
		if err != nil {
			panic(fmt.Sprintf("获取数据库连接失败: %v", err))
		}
	}
	return db
}

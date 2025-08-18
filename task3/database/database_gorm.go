package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Email string `gorm:"size:100;uniqueIndex"`
	Age   int
}

// Product 产品模型
type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Price float64
}

// 数据库连接配置
const (
	DBUser     = "root"
	DBPassword = "123456"
	DBHost     = "localhost"
	DBPort     = "3306"
	DBName     = "testdb"
)

// ConnectDB 连接数据库
func ConnectDB() (*gorm.DB, error) {
	// 构建DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功！")
	return db, nil
}

// InitDB 初始化数据库（创建表）
func InitDB(db *gorm.DB) error {
	// 自动迁移模式（创建表）
	err := db.AutoMigrate(&User{}, &Product{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Println("数据库表创建成功！")
	return nil
}

// CreateUser 创建用户
func CreateUser(db *gorm.DB, name, email string, age int) error {
	user := User{
		Name:  name,
		Email: email,
		Age:   age,
	}

	result := db.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("创建用户失败: %v", result.Error)
	}

	log.Printf("用户创建成功，ID: %d\n", user.ID)
	return nil
}

// GetUsers 获取所有用户
func GetUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("查询用户失败: %v", result.Error)
	}

	return users, nil
}

// UpdateUser 更新用户
func UpdateUser(db *gorm.DB, userID uint, name string) error {
	result := db.Model(&User{}).Where("id = ?", userID).Update("name", name)
	if result.Error != nil {
		return fmt.Errorf("更新用户失败: %v", result.Error)
	}

	log.Printf("用户更新成功，影响行数: %d\n", result.RowsAffected)
	return nil
}

// DeleteUser 删除用户
func DeleteUser(db *gorm.DB, userID uint) error {
	result := db.Delete(&User{}, userID)
	if result.Error != nil {
		return fmt.Errorf("删除用户失败: %v", result.Error)
	}

	log.Printf("用户删除成功，影响行数: %d\n", result.RowsAffected)
	return nil
}

// DatabaseExample 数据库操作示例函数
func DatabaseExample() {
	fmt.Println("开始数据库测试...")

	// 连接数据库
	db, err := ConnectDB()
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		return
	}

	// 初始化数据库（创建表）
	if err := InitDB(db); err != nil {
		log.Printf("数据库初始化失败: %v", err)
		return
	}

	// 创建用户示例
	if err := CreateUser(db, "张三", "zhangsan@example.com", 25); err != nil {
		log.Printf("创建用户错误: %v", err)
	}

	if err := CreateUser(db, "李四", "lisi@example.com", 30); err != nil {
		log.Printf("创建用户错误: %v", err)
	}

	// 查询所有用户
	users, err := GetUsers(db)
	if err != nil {
		log.Printf("查询用户错误: %v", err)
	} else {
		fmt.Println("\n所有用户:")
		for _, user := range users {
			fmt.Printf("ID: %d, 姓名: %s, 邮箱: %s, 年龄: %d\n",
				user.ID, user.Name, user.Email, user.Age)
		}
	}

	// 更新用户示例
	if len(users) > 0 {
		if err := UpdateUser(db, users[0].ID, "张三修改"); err != nil {
			log.Printf("更新用户错误: %v", err)
		}

		// 再次查询验证更新
		users, _ = GetUsers(db)
		fmt.Println("\n更新后的用户:")
		for _, user := range users {
			fmt.Printf("ID: %d, 姓名: %s, 邮箱: %s, 年龄: %d\n",
				user.ID, user.Name, user.Email, user.Age)
		}
	}

	fmt.Println("\n数据库测试完成！")
}

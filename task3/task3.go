package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL 驱动，下划线表示只导入其 side-effects
	"github.com/jmoiron/sqlx"          // 导入 sqlx 包
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQL语句练习

// --题目1：基本CRUD操作
// -- 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
// -- 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
//     insert into students(`name`,`age`,`grade`) values ('张三',20,'三年级')
// -- 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
//     select `id`,`name`,`age`,`grade` from students where `age` > 18
// -- 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
//     update students set `grade` = '四年级' where `name`='张三'
// -- 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
//     delete from students where `age` < 15

// 题目2：事务语句
// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
// 要求 ：
// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

// START TRANSACTION;
// SELECT balance FROM accounts WHERE id = A FOR UPDATE;  -- 这里的 FOR UPDATE 是什么意思
// IF (select balance FROM accounts WHERE id = A) >= 100 THEN
//     UPDATE accounts SET balance = balance-100 WHERE id = A;
//     UPDATE accounts SET balance = balance+100 WHERE id = B;
//     INSERT INTO transactions(from_account_id,to_account_id,amount) VALUES(A,B,100);
//     COMMIT;
// else
//     ROLLBACK;
// END IF;

// Sqlx入门
// 数据库配置常量
const (
	DBUser     = "root"
	DBPassword = "123456"
	DBHost     = "localhost"
	DBPort     = "3306"
	DBName     = "testdb"
)

func InitDB() (*sqlx.DB, error) {
	// 1. 构建 DSN (Data Source Name)
	// 格式通常为 "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// parseTime=True 很重要，它会将 MySQL 的 DATE/DATETIME/TIMESTAMP 类型自动解析为 Go 的 time.Time 类型
	// loc=Local 也很重要，它会使用本地时区解析时间
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)
	// 2. 连接到 MySQL 数据库
	// 使用 sqlx.Connect 连接，它会同时打开连接并尝试 ping 数据库
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		// 使用 fmt.Errorf 包装原始错误，提供更多上下文信息
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Database connection established successfully with MySQL.")
	return db, nil
}

// 题目1：使用SQL扩展库进行查询
// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
// 要求 ：
// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

// Employee 结构体用于映射 employees 表的行数据
type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func sqlxTask1() {
	// 调用 InitDB 函数来获取数据库连接
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err) // 如果连接失败，程序直接退出
	}
	defer db.Close() // 确保数据库连接在 main 函数结束时关闭

	// 3. 创建 employees 表 (如果不存在)
	// 注意 MySQL 的数据类型和 AUTO_INCREMENT 语法
	schema := `
	CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		department VARCHAR(255) NOT NULL,
		salary INT NOT NULL
	);`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to create employees table: %v", err)
	}
	fmt.Println("Employees table checked/created.")

	// 4. 插入初始测试数据 (如果表中没有数据)
	var count int
	// 使用 db.Get 获取单行单列结果
	err = db.Get(&count, "SELECT COUNT(*) FROM employees")
	if err != nil {
		log.Fatalf("Failed to get employee count: %v", err)
	}

	if count == 0 {
		insertStmt := `
		INSERT INTO employees (name, department, salary) VALUES
		('张三', '技术部', 8000),
		('李四', '销售部', 6000),
		('王五', '技术部', 9500),
		('赵六', '市场部', 7000),
		('钱七', '技术部', 10000),
		('孙八', '销售部', 6500);`
		_, err = db.Exec(insertStmt)
		if err != nil {
			log.Fatalf("Failed to insert initial data: %v", err)
		}
		fmt.Println("Initial employee data inserted.")
	} else {
		fmt.Println("Employee data already exists, skipping insertion.")
	}

	// --- 任务 1: 查询所有部门为 "技术部" 的员工信息 ---
	fmt.Println("\n--- Query 1: Employees in '技术部' ---")
	techEmployees := []Employee{} // 声明一个 Employee 结构体切片来存储结果

	// 使用 db.Select() 将查询结果直接映射到结构体切片
	// 对于 MySQL，通常使用 ? 作为占位符
	err = db.Select(&techEmployees, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")
	if err != nil {
		log.Fatalf("Failed to query tech department employees: %v", err)
	}

	if len(techEmployees) > 0 {
		for _, emp := range techEmployees {
			fmt.Printf("ID: %d, Name: %s, Dept: %s, Salary: %d\n", emp.ID, emp.Name, emp.Department, emp.Salary)
		}
	} else {
		fmt.Println("No employees found in '技术部'.")
	}

	// --- 任务 2: 查询工资最高的员工信息 ---
	fmt.Println("\n--- Query 2: Employee with Highest Salary ---")
	highestPaidEmployee := Employee{} // 声明一个 Employee 结构体来存储单个结果

	// 使用 db.Get() 将单行查询结果直接映射到结构体
	err = db.Get(&highestPaidEmployee, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		if err == sql.ErrNoRows { // 处理没有结果的情况
			fmt.Println("No employees found in the table.")
		} else {
			log.Fatalf("Failed to query highest paid employee: %v", err)
		}
	} else {
		fmt.Printf("Highest Paid Employee: ID: %d, Name: %s, Dept: %s, Salary: %d\n",
			highestPaidEmployee.ID, highestPaidEmployee.Name, highestPaidEmployee.Department, highestPaidEmployee.Salary)
	}
}

// 题目2：实现类型安全映射
// 假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
// 要求 ：
// 定义一个 Book 结构体，包含与 books 表对应的字段。
// 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

// Book 结构体用于映射 books 表的行数据
// 结构体字段的 `db:"column_name"` tag 将 SQL 列名映射到 Go 结构体字段
// 这就是 sqlx 实现类型安全映射的关键
type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"` // 使用 float64 来匹配数据库的 DECIMAL(10, 2)
}

func sqlxTask2() {
	// 调用 InitDB 函数来获取数据库连接
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err) // 如果连接失败，程序直接退出
	}
	defer db.Close() // 确保数据库连接在 main 函数结束时关闭

	// 3. 创建 books 表 (如果不存在)
	// price 使用 DECIMAL(10, 2) 存储货币更精确
	schema := `
	CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		price DECIMAL(10, 2) NOT NULL
	);`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to create books table: %v", err)
	}
	fmt.Println("Books table checked/created.")

	// 4. 插入初始测试数据 (如果表中没有数据)
	var bookCount int
	err = db.Get(&bookCount, "SELECT COUNT(*) FROM books")
	if err != nil {
		log.Fatalf("Failed to get book count: %v", err)
	}

	if bookCount == 0 {
		insertBookStmt := `
		INSERT INTO books (title, author, price) VALUES
		('Go语言编程', '许式伟', 75.50),
		('Python Cookbook', 'David Beazley', 99.00),
		('Effective Go', 'Go Team', 45.00),
		('Clean Code', 'Robert C. Martin', 60.25),
		('The Go Programming Language', 'Alan A. A. Donovan', 120.00),
		('SQL必知必会', 'Ben Forta', 38.80);`
		_, err = db.Exec(insertBookStmt)
		if err != nil {
			log.Fatalf("Failed to insert initial book data: %v", err)
		}
		fmt.Println("Initial book data inserted.")
	} else {
		fmt.Println("Book data already exists, skipping insertion.")
	}

	// --- 核心任务：查询价格大于 50 元的书籍 ---
	fmt.Println("\n--- Query: Books with Price > 50 ---")

	// 声明一个 Book 结构体切片来存储查询结果
	// 这是确保类型安全的关键：sqlx 会将每行数据映射成一个 Book 类型的实例
	expensiveBooks := []Book{}

	// 定义查询 SQL 语句
	query := "SELECT id, title, author, price FROM books WHERE price > ?"
	priceThreshold := 50.0 // 定义价格阈值，使用浮点数

	// 使用 db.Select() 执行查询并将结果映射到切片
	// sqlx 自动处理了从数据库列到结构体字段的类型转换和填充
	err = db.Select(&expensiveBooks, query, priceThreshold)
	if err != nil {
		log.Fatalf("Failed to query expensive books: %v", err)
	}

	if len(expensiveBooks) > 0 {
		fmt.Printf("Found %d books with price > %.2f:\n", len(expensiveBooks), priceThreshold)
		for _, book := range expensiveBooks {
			fmt.Printf("ID: %d, Title: \"%s\", Author: %s, Price: %.2f\n",
				book.ID, book.Title, book.Author, book.Price)
		}
	} else {
		fmt.Printf("No books found with price > %.2f.\n", priceThreshold)
	}
}

// 进阶gorm
// InitDBGORM 函数封装 GORM 数据库连接的逻辑
func InitDBGORM() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 配置 GORM 日志模式，方便调试
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database with GORM: %w", err)
	}

	fmt.Println("GORM database connection established successfully.")
	return db, nil
}

// 题目1：模型定义
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 要求 ：
// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
// 编写Go代码，使用Gorm创建这些模型对应的数据库表。

// User 模型
type User struct {
	gorm.Model          // GORM 提供的基本字段：ID, CreatedAt, UpdatedAt, DeletedAt
	Name         string `gorm:"type:varchar(100);not null;uniqueIndex"` // 用户名
	Email        string `gorm:"type:varchar(255);unique"`               // 邮箱
	Posts        []Post // User 与 Post 是一对多关系 (一个用户有多篇文章)
	PostCount    int    `gorm:"default:0"` // 用户文章数量统计字段
	CommentCount int    `gorm:"default:0"` // 用户评论数量统计字段
}

// Post 模型
type Post struct {
	gorm.Model              // GORM 提供的基本字段
	Title         string    `gorm:"type:varchar(255);not null"`     // 文章标题
	Content       string    `gorm:"type:text"`                      // 文章内容
	UserID        uint      `gorm:"not null"`                       // 外键，指向 User.ID
	User          User      `gorm:""`                               // Post 与 User 的关联
	Comments      []Comment `gorm:"foreignKey:PostID"`              // Post 与 Comment 是一对多关系
	CommentStatus string    `gorm:"type:varchar(50);default:'有评论'"` // 评论状态
	// `gorm:"default:'有评论'"` 表示默认值，GORM在创建时会设置
	CommentCount int `gorm:"default:0"` // 文章评论数量统计字段
}

// Comment 模型
type Comment struct {
	gorm.Model        // GORM 提供的基本字段
	Content    string `gorm:"type:text;not null"` // 评论内容
	UserID     uint   `gorm:"not null"`           // 外键，指向 User.ID (谁评论的)
	User       User   `gorm:""`                   // Comment 与 User 的关联
	PostID     uint   `gorm:"not null"`           // 外键，指向 Post.ID
	Post       Post   `gorm:""`                   // Comment 与 Post 的关联
}

// 题目2：关联查询
// 基于上述博客系统的模型定义。
// 要求 ：
// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
// 编写Go代码，使用Gorm查询评论数量最多的文章信息。

// 题目3：钩子函数
// 继续使用博客系统的模型。
// 要求 ：
// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

// Post 模型的 BeforeCreate 钩子 (GORM 会自动查找并执行)
// 在创建文章之前自动初始化一些字段，或者执行预处理
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	// 可以在这里设置默认值或进行验证
	if p.CommentStatus == "" {
		p.CommentStatus = "有评论" // 确保默认值被设置
	}
	return
}

// Post 模型的 AfterCreate 钩子 (GORM 会自动查找并执行)
// 在文章创建后自动更新用户的文章数量统计字段
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 使用 tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", gorm.Expr("post_count + ?", 1))
	// 这是 GORM 的原子更新语法，推荐用于计数器，避免并发问题
	err = tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", gorm.Expr("post_count + ?", 1)).Error
	if err != nil {
		log.Printf("Error updating user post_count after post creation: %v", err)
	} else {
		fmt.Printf("User %d's post_count updated after new post '%s'.\n", p.UserID, p.Title)
	}
	return
}

// Comment 模型的 BeforeDelete 钩子 (GORM 会自动查找并执行)
// 在评论删除前，更新文章的评论数量统计字段
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	// 原子递减文章评论数量
	err = tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	if err != nil {
		log.Printf("Error decrementing post comment_count before comment deletion: %v", err)
	} else {
		fmt.Printf("Post %d's comment_count decremented before comment %d deletion.\n", c.PostID, c.ID)
	}
	return
}

// Comment 模型的 AfterDelete 钩子 (GORM 会自动查找并执行)
// 在评论删除后检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 检查PostID是否有效
	if c.PostID == 0 {
		log.Printf("警告：Comment的PostID为0，跳过钩子函数")
		return nil
	}

	// 查询剩余评论数量
	var count int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		return fmt.Errorf("统计评论数量失败: %v", err)
	}

	// 准备更新数据
	updates := map[string]interface{}{
		"comment_count": count,
	}

	if count == 0 {
		updates["comment_status"] = "无评论"
	}

	// 使用明确的WHERE条件更新
	return tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(updates).Error

}

func gormTask1() {
	db, err := InitDBGORM()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// -----------------------------------------------------------------------------
	// 题目1：创建数据库表
	// -----------------------------------------------------------------------------
	fmt.Println("\n--- Creating Database Tables ---")
	// AutoMigrate 会根据模型定义自动创建表或更新表结构
	// 如果表已存在，它只会添加缺失的列，不会删除现有数据
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database tables: %v", err)
	}
	fmt.Println("Database tables created/migrated successfully.")

	// -----------------------------------------------------------------------------
	// 插入一些测试数据
	// -----------------------------------------------------------------------------
	fmt.Println("\n--- Inserting Test Data ---")

	// 插入用户
	user1 := User{Name: "Alice", Email: "alice@example.com"}
	user2 := User{Name: "Bob", Email: "bob@example.com"}

	// 检查并插入用户
	var existingUser User
	if err := db.Where("name = ?", user1.Name).First(&existingUser).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.FirstOrCreate(&user1).Error; err != nil {
			log.Fatalf("Failed to create user1: %v", err)
		}
		fmt.Printf("User created: %s (ID: %d)\n", user1.Name, user1.ID)
	} else if err == nil {
		user1 = existingUser // 使用已存在的用户
		fmt.Printf("User %s already exists (ID: %d).\n", user1.Name, user1.ID)
	} else {
		log.Fatalf("Error checking user: %v", err)
	}

	if err := db.Where("name = ?", user2.Name).First(&existingUser).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.FirstOrCreate(&user2).Error; err != nil {
			log.Fatalf("Failed to create user2: %v", err)
		}
		fmt.Printf("User created: %s (ID: %d)\n", user2.Name, user2.ID)
	} else if err == nil {
		user2 = existingUser // 使用已存在的用户
		fmt.Printf("User %s already exists (ID: %d).\n", user2.Name, user2.ID)
	} else {
		log.Fatalf("Error checking user: %v", err)
	}

	// 插入文章
	post1 := Post{Title: "Go并发编程初探", Content: "Go语言的并发特性非常强大。", UserID: user1.ID}
	post2 := Post{Title: "GORM使用指南", Content: "一个强大的Go ORM框架。", UserID: user1.ID}
	post3 := Post{Title: "区块链基础", Content: "理解区块链的核心概念。", UserID: user2.ID}

	// 检查并插入文章，触发 AfterCreate 钩子
	var existingPost Post
	if err := db.Where("title = ?", post1.Title).First(&existingPost).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&post1).Error; err != nil {
			log.Fatalf("Failed to create post1: %v", err)
		}
		fmt.Printf("Post created: \"%s\" by User %d (ID: %d)\n", post1.Title, post1.UserID, post1.ID)
	} else if err == nil {
		post1 = existingPost
		fmt.Printf("Post \"%s\" already exists (ID: %d).\n", post1.Title, post1.ID)
	} else {
		log.Fatalf("Error checking post: %v", err)
	}

	if err := db.Where("title = ?", post2.Title).First(&existingPost).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&post2).Error; err != nil {
			log.Fatalf("Failed to create post2: %v", err)
		}
		fmt.Printf("Post created: \"%s\" by User %d (ID: %d)\n", post2.Title, post2.UserID, post2.ID)
	} else if err == nil {
		post2 = existingPost
		fmt.Printf("Post \"%s\" already exists (ID: %d).\n", post2.Title, post2.ID)
	} else {
		log.Fatalf("Error checking post: %v", err)
	}

	if err := db.Where("title = ?", post3.Title).First(&existingPost).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&post3).Error; err != nil {
			log.Fatalf("Failed to create post3: %v", err)
		}
		fmt.Printf("Post created: \"%s\" by User %d (ID: %d)\n", post3.Title, post3.UserID, post3.ID)
	} else if err == nil {
		post3 = existingPost
		fmt.Printf("Post \"%s\" already exists (ID: %d).\n", post3.Title, post3.ID)
	} else {
		log.Fatalf("Error checking post: %v", err)
	}

	// 插入评论
	comment1 := Comment{Content: "文章写得真棒！", UserID: user2.ID, PostID: post1.ID}
	comment2 := Comment{Content: "非常实用，感谢分享！", UserID: user1.ID, PostID: post1.ID}
	comment3 := Comment{Content: "GORM简化了开发。", UserID: user2.ID, PostID: post2.ID}
	comment4 := Comment{Content: "期待更多相关文章。", UserID: user1.ID, PostID: post3.ID}

	// 检查并插入评论
	var existingComment Comment
	commentsToCreate := []Comment{comment1, comment2, comment3, comment4}
	for i, c := range commentsToCreate {
		if err := db.Where("post_id = ? AND user_id = ? AND content = ?", c.PostID, c.UserID, c.Content).First(&existingComment).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&commentsToCreate[i]).Error; err != nil {
				log.Fatalf("Failed to create comment: %v", err)
			}
			// 同时更新 Post 的 CommentCount
			db.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
			// 同时更新 User 的 CommentCount
			db.Model(&User{}).Where("id = ?", c.UserID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
			fmt.Printf("Comment created for Post %d by User %d: \"%s\" (ID: %d)\n", c.PostID, c.UserID, c.Content, commentsToCreate[i].ID)
		} else if err == nil {
			commentsToCreate[i] = existingComment
			fmt.Printf("Comment \"%s\" for Post %d already exists (ID: %d).\n", c.Content, c.PostID, commentsToCreate[i].ID)
		} else {
			log.Fatalf("Error checking comment: %v", err)
		}
	}

	// -----------------------------------------------------------------------------
	// 题目2：关联查询
	// -----------------------------------------------------------------------------
	fmt.Println("\n--- Associative Queries ---")

	// 2.1 查询某个用户（例如 Alice）发布的所有文章及其对应的评论信息
	fmt.Println("\n--- Query: Alice's Posts and their Comments ---")
	var aliceUser User
	// 使用 Preload("Posts.Comments") 预加载关联的 Posts 和 Comments
	// `Preload("Posts")` 预加载用户的文章，`Preload("Posts.Comments")` 预加载文章的评论
	// `Preload("Comments.User")` 可以预加载评论的用户
	err = db.Where("name = ?", "Alice").Preload("Posts.Comments").First(&aliceUser).Error
	if err != nil {
		log.Fatalf("Failed to query Alice's posts with comments: %v", err)
	}

	fmt.Printf("User: %s (ID: %d, PostCount: %d, CommentCount: %d)\n", aliceUser.Name, aliceUser.ID, aliceUser.PostCount, aliceUser.CommentCount)
	for _, post := range aliceUser.Posts {
		fmt.Printf("  Post: \"%s\" (ID: %d, CommentCount: %d, Status: %s)\n", post.Title, post.ID, post.CommentCount, post.CommentStatus)
		for _, comment := range post.Comments {
			fmt.Printf("    Comment ID: %d, Content: \"%s\"\n", comment.ID, comment.Content)
		}
	}

	// 2.2 查询评论数量最多的文章信息
	fmt.Println("\n--- Query: Post with Most Comments ---")
	var mostCommentedPost Post
	// 使用 Select("*, (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count")
	// 这里通过子查询计算评论数量，并作为虚拟列返回。
	// GORM 也可以直接使用 Count 方法配合关联查询，但这里直接聚合更直接。
	// Note: 这种直接在 SELECT 中计算的方式，`Post.CommentCount` 需要从数据库加载，而不是 GORM 自动更新。
	// GORM 推荐的方法是在模型定义中使用 `gorm:"-"` 忽略掉 CommentCount 字段，然后通过 `db.Joins()` 或 `db.Raw()` 自定义聚合查询。
	// 为了简便，这里直接使用子查询或者直接查询 Post 表的 CommentCount 字段（假设已被钩子更新）。
	// 如果 CommentCount 是由钩子维护的，直接 ORDER BY 即可
	err = db.Order("comment_count DESC").First(&mostCommentedPost).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("No posts found.")
		} else {
			log.Fatalf("Failed to query post with most comments: %v", err)
		}
	} else {
		fmt.Printf("Post with most comments: \"%s\" (ID: %d, Comments: %d)\n",
			mostCommentedPost.Title, mostCommentedPost.ID, mostCommentedPost.CommentCount)
	}

	// -----------------------------------------------------------------------------
	// 题目3：钩子函数测试 (Post 的 AfterCreate 已经在上面插入数据时触发了)
	// 现在测试 Comment 的 BeforeDelete 和 AfterDelete 钩子
	// -----------------------------------------------------------------------------
	fmt.Println("\n--- Testing Comment Delete Hooks ---")

	// 查找 post1 的一个评论
	var commentToDelete Comment
	err = db.Where("post_id = ? AND content = ?", post1.ID, "文章写得真棒！").First(&commentToDelete).Error
	if err != nil {
		log.Fatalf("Failed to find comment to delete: %v", err)
	}
	fmt.Printf("Attempting to delete comment ID %d from Post ID %d.\n", commentToDelete.ID, commentToDelete.PostID)

	// 删除评论，触发 BeforeDelete 和 AfterDelete 钩子
	err = db.Delete(&commentToDelete).Error
	if err != nil {
		log.Fatalf("Failed to delete comment: %v", err)
	}
	fmt.Printf("Comment ID %d deleted successfully.\n", commentToDelete.ID)

	// 再次删除 post1 的另一个评论
	var anotherCommentToDelete Comment
	err = db.Where("post_id = ? AND content = ?", post1.ID, "非常实用，感谢分享！").First(&anotherCommentToDelete).Error
	if err != nil {
		log.Fatalf("Failed to find another comment to delete: %v", err)
	}
	fmt.Printf("Attempting to delete comment ID %d from Post ID %d.\n", anotherCommentToDelete.ID, anotherCommentToDelete.PostID)

	err = db.Where("id = ?", anotherCommentToDelete.ID).Delete(&Comment{}).Error
	if err != nil {
		log.Fatalf("Failed to delete another comment: %v", err)
	}
	fmt.Printf("Comment ID %d deleted successfully.\n", anotherCommentToDelete.ID)

	// 验证文章的评论数量和状态 (重新加载文章)
	var updatedPost1 Post
	err = db.First(&updatedPost1, post1.ID).Error
	if err != nil {
		log.Fatalf("Failed to get updated post1: %v", err)
	}
	fmt.Printf("\nUpdated Post 1 Status: Title=\"%s\", CommentCount: %d, CommentStatus: %s\n",
		updatedPost1.Title, updatedPost1.CommentCount, updatedPost1.CommentStatus)

	if updatedPost1.CommentCount == 0 && updatedPost1.CommentStatus == "无评论" {
		fmt.Println("Comment deletion hooks worked as expected for Post 1!")
	} else {
		fmt.Println("Comment deletion hooks did not update Post 1 as expected.")
	}

	// 验证用户文章计数
	var updatedUser1 User
	db.First(&updatedUser1, user1.ID)
	fmt.Printf("\nUpdated User %s PostCount: %d\n", updatedUser1.Name, updatedUser1.PostCount)
	var updatedUser2 User
	db.First(&updatedUser2, user2.ID)
	fmt.Printf("Updated User %s PostCount: %d\n", updatedUser2.Name, updatedUser2.PostCount)
}

func main() {
	// sqlxTask1()
	// sqlxTask2()
	gormTask1()
}

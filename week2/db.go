package main

import (
	"context"
	"encoding/json"

	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:jack_040604@tcp(127.0.0.1:3306)/db01?charset=utf8"

var db *gorm.DB
var rdb *redis.Client

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 错误处理封装
func HandleError(c *gin.Context, err error, message string) {
	log.Printf("%s: %v", message, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

// Redis缓存处理(读)
func RedisCacheGetId(c *gin.Context) (*Student, error) {
	id := c.Param("id")
	// 先查看Json里面有无数据
	studentbt, err := rdb.Get(c.Request.Context(), id).Bytes()
	if err != nil {
		if err != redis.Nil {
			// 如果发生了除缓存不存在以外的其他错误，则返回500错误
			HandleError(c, err, "Failed to retrieve data from cache")
			return nil, fmt.Errorf("failed to retrieve data from cache")
		}
		// 到MySQL中查找
		student, err := FindInMysql(id, c)
		// 未找到则返回错误
		if err != nil {
			return nil, err
		}
		if student == nil {
			return nil, nil
		}
		// 缓存到Redis中
		if _, err := AddToRedisCache(student, id); err != nil {
			return nil, err
		}
		return student, nil
	}
	var student Student
	err = json.Unmarshal(studentbt, &student)
	if err != nil {
		return nil, err
	}
	return &student, nil

}

// 缓存Redis
func AddToRedisCache(student *Student, id string) (*Student, error) {
	studentJSON, err := json.Marshal(student)
	if err != nil {
		return student, err
	}
	if err := rdb.Set(context.Background(), id, studentJSON, time.Hour).Err(); err != nil {
		HandleError(nil, err, "Failed to cache data in Redis")
		return nil, err
	}
	return student, nil
}

// 删除Redis中的缓存
func DeletRedis(c *gin.Context) error {
	err := rdb.Del(context.Background(), c.Param("id")).Err()
	if err != nil {
		// if err := SendMessage(c.Param("id")); err != nil {
		// 	// 如果发送消息失败，则返回错误，但是仍然尝试删除缓存
		// 	log.Printf("Failed to send message to queue: %v", err)
		// }
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to Delete data in Redis"})
		return err
	}
	return nil
}

// 在MySQL中查找数据
func FindInMysql(id string, c *gin.Context) (*Student, error) {
	var student *Student
	studentID, _ := strconv.Atoi(id)
	// fmt.Println(id)
	if err := db.First(&student, studentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		HandleError(c, err, "Database query error")
		return nil, err
	}
	return student, nil
}

func InitDB() error {
	//使用gorm建立Mysql连接池
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL:%v", err)
	}
	Mysqldb, _ := db.DB()
	if err := Mysqldb.Ping(); err != nil {
		log.Fatalf("Failed to Ping MySQL:%v", err)
	}
	fmt.Println("MySQL sucessfully connect")
	Mysqldb.SetMaxOpenConns(64)
	Mysqldb.SetMaxIdleConns(64)
	Mysqldb.SetConnMaxLifetime(5 * time.Minute)
	//初始化Redis连接池
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// defer rdb.Close()
	Pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to ping Redis:%v", err)
	}
	fmt.Println("Redis Ping Response:", Pong)
	return nil
}

// 判断id
func JudgeId(c *gin.Context) int {
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid student ID"})
		return 0
	}
	return studentID
}
func GetAllStudents() ([]Student, error) {
	var students []Student
	if err := db.Find(&students).Error; err != nil {
		// Find(&students) 用于从数据库中检索所有的学生信息，并将结果存储到 students 变量中
		return nil, err
	}
	// 将所有学生数据缓存到Redis缓存中
	for _, student := range students {
		_, err := AddToRedisCache(&student, strconv.Itoa(student.ID))
		if err != nil {
			return nil, err
		}
	}
	return students, nil
}

func ListAllStudents(c *gin.Context) {
	students, err := GetAllStudents()
	if err != nil {
		HandleError(c, err, "Failed to get all students")
		return
	}

	c.JSON(http.StatusOK, students)
}

func AddStudent(c *gin.Context) {
	var student Student
	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&student).Error; err != nil {
		HandleError(c, err, "Failed to add student to database")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sucessfully added"})
}

// GetStudentByID 获取单个学生信息
func GetStudentById(c *gin.Context) {
	student, err := RedisCacheGetId(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No this student"})
		return
	}
	c.JSON(http.StatusOK, student)
}

func UpdataStudent(c *gin.Context) {
	var student Student
	studentID := JudgeId(c)
	if studentID == 0 {
		return
	}

	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind JSON"})
		return
	}
	//更新MySQL中的数据
	if err := db.Model(&student).Where("id = ?", studentID).Updates(Student{Name: student.Name, Age: student.Age}).Error; err != nil {
		HandleError(c, err, "Failed to update student in database")
		return
	}
	//删去缓存
	if err := DeletRedis(c); err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully updated"})
}

func DeleteStudent(c *gin.Context) {
	studentID := JudgeId(c)
	if studentID == 0 {
		return
	}

	var student Student
	if err := db.First(&student, studentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		return
	}

	c.Set("student", student) // 将学生信息存储到 Gin 上下文中
	student = c.MustGet("student").(Student)
	if err := db.Delete(&student).Error; err != nil {
		HandleError(c, err, "Failed to delete student")
		return
	}
	//删去缓存
	if err := DeletRedis(c); err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

// 测试
func GetStudentDB(c *gin.Context) {
	id := c.Param("id")
	student, err := FindInMysql(id, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到学生"})
	} else {
		c.JSON(http.StatusOK, student)
	}
}

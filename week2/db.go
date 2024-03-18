package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var dsn = "root:jack_040604@tcp(127.0.0.1:3306)/db01?charset=utf8"
var Mysqldb *sql.DB

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func InitDB() error {
	var err error
	//初始化Mysql连接池
	Mysqldb, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect MySQL:%v", err)
	}

	if err := Mysqldb.Ping(); err != nil {
		log.Fatalf("Failed to Ping MySQL:%v", err)
	}

	return nil
}
func GetAllStudents() ([]Student, error) {
	var students []Student
	// Mysqldb.Query("SELECT id, name, age FROM students") 这行代码执行了一个 SELECT 查询，从名为 "students" 的表中检索所有学生的 ID、姓名和年龄信息。查询的结果被存储在 rows 变量中，以供后续使用。
	// if err := Mysqldb.Find(&students).Error(); err != nil {
	// 	return nil, err
	// }
	rows, err := Mysqldb.Query("SELECT id, name, age FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 当所有行都被遍历完成后，您需要调用 rows.Close() 来关闭结果集

	// 迭代查询结果集中的每一行
	for rows.Next() {
		var student Student
		// rows.Scan() 方法会根据传入的参数列表，依次从当前行中读取每个列的值，并将这些值分别赋给相应的变量。
		if err := rows.Scan(&student.ID, &student.Name, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func ListAllStudents(c *gin.Context) {
	students, err := GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	//  SQL 插入语句，将学生的姓名和年龄插入到名为 students 的表中(序号更新问题需要到这里来解决)
	_, err := Mysqldb.Exec("INSERT INTO students (name, age) VALUES (?, ?)", student.Name, student.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sucessfully added"})
}

// GetStudentByID 获取单个学生信息
func GetStudentById(c *gin.Context) {
	// studentid := c.Param("id")
	// c.Param("id") 用于获取 URL 中的参数值。
	// id := strconv.Atoi()
	// 函数正是用于将字符串转换为整数的。
	var student Student
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No find"})
	}
	row := Mysqldb.QueryRow("SELECT id, name, age FROM students WHERE id = ?", studentID)
	if err := row.Scan(&student.ID, &student.Name, &student.Age); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, student)
}

func UpdataStudent(c *gin.Context) {
	var student Student
	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "wrong updata"})
		return
	}
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid student ID"})
		return
	}
	_, err = Mysqldb.Exec("UPDATE students SET name=?, age=? WHERE id=?", student.Name, student.Age, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully updated"})
}

func DeleteStudent(c *gin.Context) {
	// student := c.MustGet("student").(Student)
	// if err := Mysqldb.Delete(&student).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	studentID, _ := strconv.Atoi(c.Param("id"))
	_, err := Mysqldb.Exec("DELETE FROM students WHERE id=?", studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// API路由组
	api := router.Group("/api")

	// 使用中间件检查学生是否存在
	api.Use(CheckStudentExist())

	// 获取所有学生信息
	api.GET("/students", ListStudents)

	// 创建学生信息
	api.Use(CheckStudentExist()).POST("/students", addStudent)

	// 获取单个学生信息
	api.Use(CheckStudentExist()).GET("/students/:id", getStudentByID)

	// 更新学生信息
	api.Use(CheckStudentExist()).PUT("/students/:id", updateStudent)

	// 删除学生信息
	api.Use(CheckStudentExist()).DELETE("/students/:id", deleteStudent)

	return router
}

// CheckStudentExist 中间件用于检查学生是否存在
func CheckStudentExist() gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("id")
		id, err := strconv.Atoi(studentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
			c.Abort()
			return
		}

		student := GetStudentByID(id)
		if student.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生未找到"})
			c.Abort()
			return
		}

		// 将学生信息添加到上下文中，以便后续处理函数可以直接获取
		c.Set("student", student)
		c.Next()
	}
}

// ListStudents 返回所有学生信息
func ListStudents(c *gin.Context) {
	students, err := GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

// AddStudent 创建学生信息
func addStudent(c *gin.Context) {
	var student Student
	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	AddStudent(student)
	c.JSON(http.StatusOK, gin.H{"message": "学生信息添加成功"})
}

// GetStudentByID 获取单个学生信息
func getStudentByID(c *gin.Context) {
	student := c.MustGet("student").(Student)
	c.JSON(http.StatusOK, student)
}

// UpdateStudent 更新学生信息
func updateStudent(c *gin.Context) {
	var student Student
	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	studentToUpdate := c.MustGet("student").(Student)
	student.ID = studentToUpdate.ID
	UpdateStudent(student)
	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}

// DeleteStudent 删除学生信息
func deleteStudent(c *gin.Context) {
	student := c.MustGet("student").(Student)
	DeleteStudent(student.ID)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

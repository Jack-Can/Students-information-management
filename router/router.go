package router

import (
	"net/http"
	"strconv"
	"student/dbservice"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// API路由组
	api := router.Group("/api")

	// 获取所有学生信息
	api.GET("/students", ListStudents)

	// 创建学生信息
	api.POST("/students", addStudent)

	// 获取单个学生信息
	api.GET("/students/:id", getStudentByID)

	// 更新学生信息
	api.PUT("/students/:id", updateStudent)

	// 删除学生信息
	api.DELETE("/students/:id", deleteStudent)

	return router
}

// ListStudents 返回所有学生信息
func ListStudents(c *gin.Context) {
	students, err := dbservice.GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

// CreateStudent  创建学生信息
func addStudent(c *gin.Context) {
	var student dbservice.Student
	if err := c.BindJSON(&student); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	dbservice.AddStudent(student)
	c.JSON(200, gin.H{"message": "学生信息添加成功"})
}

// GetStudent  获取单个学生信息
func getStudentByID(c *gin.Context) {
	studentID := c.Param("id")
	id, _ := strconv.Atoi(studentID)
	student := dbservice.GetStudentByID(id)
	if student.ID == 0 {
		c.JSON(404, gin.H{"error": "学生未找到"})
	} else {
		c.JSON(200, student)
	}
}

// UpdateStudent  更新学生信息
func updateStudent(c *gin.Context) {
	var student dbservice.Student
	if err := c.BindJSON(&student); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求体"})
		return
	}
	studentID := c.Param("id")
	id, _ := strconv.Atoi(studentID)
	hasStudent := dbservice.GetStudentByID(id)
	if hasStudent.ID == 0 {
		c.JSON(404, gin.H{"error": "学生未找到"})
	} else {
		student.ID = id
		dbservice.UpdateStudent(student)
		c.JSON(200, gin.H{"message": "修改成功"})
	}

}

// DeleteStudent  删除学生信息
func deleteStudent(c *gin.Context) {
	studentID := c.Param("id")
	id, _ := strconv.Atoi(studentID)
	student := dbservice.GetStudentByID(id)
	if student.ID == 0 {
		c.JSON(404, gin.H{"error": "学生未找到"})
	} else {
		dbservice.DeleteStudent(id)
		c.JSON(200, gin.H{"message": "删除成功"})
	}
}

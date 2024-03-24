package main

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()
	apiGroup := router.Group("/api")

	// 获得所有学生
	apiGroup.GET("/students", ListAllStudents)

	//创建学生信息
	apiGroup.POST("/student", AddStudent)
	// 获得单个学生信息
	apiGroup.GET("/student/:id", GetStudentById)
	// :id为参数占位符

	// 更新学生信息
	apiGroup.PUT("/student/:id", UpdataStudent)
	// 删除学生信息
	apiGroup.DELETE("/student/:id", DeleteStudent)
	// 对照组
	apiGroup.GET("/dbstudent/:id", GetStudentDB)
	return router
}

package dbservice

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// Student 结构表示学生信息
type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func InitDB() error {
	var err error
	db, err = gorm.Open("mysql", "root:jack_040604@/db01?charset=utf8")
	if err != nil {
		return err
	}

	db.Table("students").AutoMigrate(&Student{})
	return nil
}

func CloseDB() {
	db.Close()
	// 在这里关闭数据库连接
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GetAllStudents 获取所有学生信息
func GetAllStudents() ([]Student, error) {
	var students []Student
	if err := db.Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

// CreateStudent 创建一个学生信息
func AddStudent(students Student) {
	db, err := sql.Open("mysql", "root:jack_040604@/db01?charset=utf8")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("INSERT INTO students (name, age) VALUES (?, ?)", students.Name, students.Age)
	checkErr(err)
}

// GetStudentByID 获取单个学生信息
func GetStudentByID(studentID int) (students Student) {
	db, err := sql.Open("mysql", "root:jack_040604@/db01?charset=utf8")
	checkErr(err)
	defer db.Close()
	db.QueryRow("SELECT id, name, age FROM students WHERE id = ?", studentID).Scan(&students.ID, &students.Name, &students.Age)
	checkErr(err)
	return
}

// UpdateStudent 更新学生信息
func UpdateStudent(students Student) {
	db, err := sql.Open("mysql", "root:jack_040604@/db01?charset=utf8")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("UPDATE students SET name = ?, age = ? WHERE id = ?", students.Name, students.Age, students.ID)
	checkErr(err)
}

// DeleteStudent 删除学生信息
func DeleteStudent(id int) error {
	return db.Delete(&Student{}, id).Error
}

package main

import (
	"student/dbservice"
	"student/router"
)

func main() {
	if err := dbservice.InitDB(); err != nil {
		panic("Failed to connect to the database")
	}
	defer dbservice.CloseDB()

	router := router.SetupRouter()
	router.Run(":8080")
}

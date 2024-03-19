package main

func main() {
	if err := InitDB(); err != nil {
		panic("Filed to connect to database")
	}
	// defer InitDB()
	router := SetupRouter()
	router.Run(":8080")
}

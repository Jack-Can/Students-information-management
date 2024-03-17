package main

func main() {
	if err := InitDB(); err != nil {
		panic("Failed to connect to the database")
	}
	defer CloseDB()

	router := SetupRouter()
	router.Run(":8080")
}

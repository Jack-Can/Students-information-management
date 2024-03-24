package main

import (
	"os"
	"runtime/pprof"
)

func main() {
	if err := InitDB(); err != nil {
		panic("Filed to connect to database")
	}
	// defer InitDB()
	// go ConsumeMessages()

	// 性能分析
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	router := SetupRouter()
	router.Run(":8080")
}

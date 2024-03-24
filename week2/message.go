package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func SendMessage(id string) error {
	fmt.Println("Sending message to queue:", id)
	return nil
}

func GetMessage() (string, error) {
	// 在这里实现从消息队列获取消息的逻辑
	return "id", nil
}
func ConsumeMessages() {
	for {
		// 从消息队列中获取需要删除的键
		id, err := GetMessage()
		if err != nil {
			log.Println("Error consuming message:", err)
			continue
		}
		// 尝试删除缓存直到成功
		for {
			err := rdb.Del(context.Background(), id).Err()
			if err == nil {
				fmt.Println("Successfully deleted key from cache:", id)
				break
			}
			fmt.Println("Error deleting key from cache:", err)
			time.Sleep(time.Second) // 可以根据需求调整重试间隔
		}
	}
}

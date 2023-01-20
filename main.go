package main

import (
	"fmt"
	"node/fs"
	"node/listener"
	"node/service"
	"time"
)

func main() {
	fmt.Println("Node started!")

	storage := fs.New("/home/uroborosq/Рабочий стол/Одиночные проекты/go-dfs/Полигон/node1")
	s := service.New(storage)
	listen := listener.CreateTcpListener("12345", s)
	for {
		err := listen.Listen()
		if err != nil {
			fmt.Printf("Error occured while working! Error description: %s\n", err.Error())
			fmt.Println("Restarting in 5 seconds...")
		}
		time.Sleep(5 * time.Second)
	}
}

package main

import (
	"fmt"
	"log"
	"time"

	"ip_mole/config"
	"ip_mole/service"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	ticker := time.NewTicker(10 * time.Minute)
	for {
		fmt.Printf("检测开始......\n")
		service.CheckIPs(*cfg)
		<-ticker.C
	}
}

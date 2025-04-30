package main

import (
	"fmt"
	"log"
	"time"

	"ip_pool_chick/config"
	"ip_pool_chick/dao"
	"ip_pool_chick/service"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化 DB
	if err := dao.InitDB(cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 3. 还原is_top标记
	go service.RestTopFlag(dao.DB, *cfg)

	// 4. 初始化 Redis
	dao.InitRedis(cfg.Redis)

	// 5. 启动验证与自动失效循环
	go service.ValidatorLoop(dao.DB, *cfg)

	// 6. 简易 refiller：定时检查并补充至 20
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			length, _ := dao.Rdb.LLen(dao.Ctx, cfg.Redis.PoolKey).Result()
			if length < cfg.Pool.IpNumber {
				if ip, err := service.GetNewIPFromDB(dao.DB); err == nil {
					fmt.Printf("添加IP {%s}\n", ip)
					service.PushNewIP(ip, cfg.Redis)
				}
			}
		}
	}()

	select {} // 阻塞主 goroutine
}

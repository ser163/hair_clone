package main

import (
	"ipghost/config"
	"ipghost/lib"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	conf config.Config
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig(&conf); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	ticker := time.NewTicker(conf.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		proxies, err := lib.GetProxies(conf)
		if err != nil {
			log.Printf("获取代理失败: %v", err)
			continue
		}

		var wg sync.WaitGroup
		for _, proxy := range proxies {
			wg.Add(1)
			go func(p string, conf config.Config) {
				defer wg.Done()
				lib.CheckAndSaveProxy(p, conf)
			}(proxy, conf)
		}
		wg.Wait()
	}
}

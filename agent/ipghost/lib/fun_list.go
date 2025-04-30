package lib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"ipghost/config"
	"ipghost/db"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ProxyResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Proxies []string `json:"proxies"`
		Count   int      `json:"count"`
	} `json:"data"`
}

func GetProxies(conf config.Config) ([]string, error) {
	resp, err := http.Get(conf.ProxyAPI)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回非200状态码: %d", resp.StatusCode)
	}

	var response ProxyResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("API返回错误: %s", response.Message)
	}

	return response.Data.Proxies, nil
}

func CheckAndSaveProxy(proxyAddr string, conf config.Config) {
	isAvailable := checkProxy(proxyAddr, conf)

	host, portStr, err := net.SplitHostPort(proxyAddr)
	if err != nil {
		log.Printf("代理地址解析失败: %s, %v", proxyAddr, err)
		return
	}

	// 将端口转换为整数
	var port int
	_, err = fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		log.Printf("端口解析失败: %s, %v", portStr, err)
		return
	}

	now := time.Now()

	if isAvailable {
		// 检查代理是否已存在
		var id int64
		var checkNum, succ, fail int
		var score int

		// 初始化数据库连接
		var dbCon *sql.DB

		dbCon, err := db.InitDB(conf)
		if err != nil {
			log.Fatalf("连接数据库失败: %v", err)
		}

		row := dbCon.QueryRow("SELECT id, check_num, succ, fail, score FROM checklist WHERE ip_addr = ? AND port = ?", host, port)
		err = row.Scan(&id, &checkNum, &succ, &fail, &score)

		if err == sql.ErrNoRows {
			// 新增代理记录
			status := 1
			succ = 1

			score = calculateScore(1, succ, fail)

			agentName := conf.AgentName
			_, err = dbCon.Exec(
				"INSERT INTO checklist (ip_addr, port, protocol, is_available, last_check, check_num, succ, fail, created_time, update_time, score, agent) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				host, port, "http", status, now, 1, succ, fail, now, now, score, agentName,
			)
			if err != nil {
				log.Printf("保存代理失败: %v", err)
			} else {
				log.Printf("新增代理: %s, 可用状态: %v", proxyAddr, isAvailable)
			}
		} else if err == nil {
			// 更新现有代理记录
			checkNum++
			status := 1
			succ++

			score = calculateScore(checkNum, succ, fail)

			_, err = dbCon.Exec(
				"UPDATE checklist SET is_available = ?, last_check = ?, check_num = ?, succ = ?, fail = ?, update_time = ?, score = ? WHERE id = ?",
				status, now, checkNum, succ, fail, now, score, id,
			)
			if err != nil {
				log.Printf("更新代理失败: %v", err)
			} else {
				log.Printf("更新代理: %s, 可用状态: %v, 评分: %d", proxyAddr, isAvailable, score)
			}
		} else {
			log.Printf("查询代理记录失败: %v", err)
		}

		dbCon.Close()
	}
}

// 计算代理评分
func calculateScore(checkNum, succ, fail int) int {
	if checkNum == 0 {
		return 0
	}
	// 简单的评分算法：成功率 * 100
	return (succ * 100) / checkNum
}

func checkProxy(proxyAddr string, conf config.Config) bool {
	// 步骤1: 端口检测
	if !portCheck(proxyAddr, conf) {
		//log.Printf("代理 %s 端口检测失败", proxyAddr)
		return false
	}

	// 步骤2: HTTP访问检测
	if !httpCheck(proxyAddr, conf) {
		//log.Printf("代理 %s HTTP检测失败", proxyAddr)
		return false
	}

	log.Printf("代理 %s 检测通过", proxyAddr)
	return true
}

func portCheck(proxyAddr string, conf config.Config) bool {
	conn, err := net.DialTimeout("tcp", proxyAddr, conf.PortTimeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func httpCheck(proxyAddr string, conf config.Config) bool {
	proxyURL, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		log.Printf("无效代理地址: %s", proxyAddr)
		return false
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: conf.HTTPTimeout,
	}

	resp, err := client.Get(conf.TargetURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

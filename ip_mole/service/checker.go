package service

import (
	"fmt"
	"ip_mole/config"
	"ip_mole/dao"
	"log"
	"net"
	"net/http"
	"net/url"
)

func checkPort(ip string, conf config.CheckConfig) bool {
	conn, err := net.DialTimeout("tcp", ip, conf.PortTimeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func checkHTTP(ip string, conf config.CheckConfig) bool {
	proxyURL, err := url.Parse("http://" + ip)
	if err != nil {
		return false
	}
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		Timeout:   conf.HTTPTimeout,
	}
	resp, err := client.Get(conf.TargetURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func CheckIPs(conf config.Config) {
	db, err := dao.InitDB(conf)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
		return
	}

	var checkConfig = conf.Check
	rows, err := db.Query("SELECT id, ip_addr, port, fail, succ FROM checklist")
	if err != nil {
		fmt.Println("查询失败:", err)
		return
	}

	defer rows.Close()
	defer db.Close()

	for rows.Next() {
		var id int
		var ip string
		var port int
		var fail int
		var succ int
		if err := rows.Scan(&id, &ip, &port, &fail, &succ); err != nil {
			continue
		}
		full := fmt.Sprintf("%s:%d", ip, port)

		portOk := checkPort(full, checkConfig)
		httpOk := checkHTTP(full, checkConfig)

		if portOk && httpOk {
			if succ > 24 {
				_, err = db.Exec(`
                UPDATE checklist 
                SET last_check=NOW(), succ=1, check_num=1, fail=0, score=0, is_available=1
                WHERE id=?`, id)
				if err != nil {
					fmt.Println("更新成功状态失败:", err)
				}
			} else {
				// 检测成功
				_, err = db.Exec(`
                UPDATE checklist 
                SET last_check=NOW(), succ=succ+1, check_num=check_num+1, is_available=1 
                WHERE id=?`, id)
				if err != nil {
					fmt.Println("更新成功状态失败:", err)
				}
			}
		} else {
			// 检测失败，增加 fail 次数
			fail += 1
			if fail >= 3 {
				// 失败次数超过3次，直接删除
				_, err = db.Exec("DELETE FROM checklist WHERE id=?", id)
				if err != nil {
					fmt.Println("删除失败:", err)
				} else {
					fmt.Printf("IP %s 失败超过3次，已删除\n", full)
				}
			} else {
				// 失败次数未超过3次，继续更新
				_, err = db.Exec(`
                    UPDATE checklist 
                    SET last_check=NOW(), fail=fail+1, check_num=check_num+1 
                    WHERE id=?`, id)
				if err != nil {
					fmt.Println("更新失败状态失败:", err)
				}
			}
		}
	}
}

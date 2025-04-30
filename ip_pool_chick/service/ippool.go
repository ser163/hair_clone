package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"ip_pool_chick/config"
	"ip_pool_chick/dao"
)

// PushNewIP 将新 IP 推入 Redis 池中（去重）
func PushNewIP(ip string, redisCfg config.RedisConfig) error {
	if exists, _ := dao.Rdb.SIsMember(dao.Ctx, redisCfg.SetKey, ip).Result(); exists {
		return nil
	}
	dao.Rdb.LPush(dao.Ctx, redisCfg.PoolKey, ip)
	dao.Rdb.SAdd(dao.Ctx, redisCfg.SetKey, ip)
	return nil
}

// PopIP 从 Redis 池出栈一个 IP
func PopIP(redisCfg config.RedisConfig) (string, error) {
	ip, err := dao.Rdb.RPop(dao.Ctx, redisCfg.PoolKey).Result()
	if err != nil {
		return "", err
	}
	dao.Rdb.SRem(dao.Ctx, redisCfg.SetKey, ip)
	return ip, nil
}

// GetNewIPFromDB 从 MySQL 筛选一个新 IP
func GetNewIPFromDB(db *sql.DB) (string, error) {
	row := db.QueryRow(
		`SELECT id, ip_addr, port, protocol
         FROM checklist
         WHERE is_available=1 AND fail = 0 AND is_top= 0
         ORDER BY score DESC, succ DESC
         LIMIT 1`)
	var ip, protocol string
	var port int
	var id int
	fmt.Sprintf("找到合适IP %s \n", ip)
	if err := row.Scan(&id, &ip, &port, &protocol); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		fmt.Println(err)
		return "", err
	}
	_, err := db.Exec(`
                UPDATE checklist 
                SET is_top = 1 
                WHERE id=?`, id)
	if err != nil {
		fmt.Println("更新成功状态失败:", err)
	}

	return fmt.Sprintf("%s:%d", ip, port), nil
}

// PortCheck 检测端口连通性
func PortCheck(proxyAddr string, conf config.CheckConfig) bool {
	conn, err := net.DialTimeout("tcp", proxyAddr, conf.PortTimeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// HTTPCheck 检测 HTTP 代理可用性
func HTTPCheck(proxyAddr string, conf config.CheckConfig) bool {
	proxyURL, err := url.Parse("http://" + proxyAddr)
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

// ValidatorLoop 定时检测 Redis 中所有 IP，可用则保留，否则剔除
func ValidatorLoop(db *sql.DB, conf config.Config) {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		ips, _ := dao.Rdb.LRange(dao.Ctx, conf.Redis.PoolKey, 0, -1).Result()
		for _, ip := range ips {
			if !PortCheck(ip, conf.Check) || !HTTPCheck(ip, conf.Check) {
				dao.Rdb.LRem(dao.Ctx, conf.Redis.PoolKey, 0, ip)
				dao.Rdb.SRem(dao.Ctx, conf.Redis.SetKey, ip)

				_, err := db.Exec(`
                UPDATE checklist 
                SET last_check=NOW(), succ=0, check_num=0, fail = 0, is_top = 0, is_available=0, score=0
                WHERE ip_addr=?`, ip)
				if err != nil {
					fmt.Println("更新成功状态失败:", err)
				}
			}
		}
	}
}

func RestTopFlag(db *sql.DB, conf config.Config) {
	_, err := db.Exec(`
                UPDATE checklist 
                SET is_top = 0`)
	if err != nil {
		fmt.Println("更新成功状态失败:", err)
	}
	fmt.Println("还原is_top 标记")
}

# HAIR_CLONE · 代理工作池系统

![HAIR_CLONE Logo](./logo.png)

> 灵感来自《西游记》中孙悟空的“毫毛分身”法术 —— 一个轻量、模块化、可扩展的代理 IP 管理系统。

一个为代理 IP 池打造的多模块系统，支持收集、验证、管理功能，使用 Go 多模块工作区组织，易扩展、易维护。

HAIR_CLONE 将代理池拆分为多个灵活组件（如毫毛化身），包括：
- 🧲 多源 IP 收集器（Agent）
- 🧪 高效 IP 检测器（Mole）
- 🧰 实时堆栈池管理器（Pool Chick）

---

## 📦 项目结构

```
HAIR_CLONE/
├─ agent/ipghost/       # IP 收集客户端（支持多种代理源）
├─ ip_pool_chick/       # IP 存储池，基于 Redis 堆栈管理
├─ ip_mole/             # IP 检测器，定期验证 IP 可用性并维护数据库
├─ go.work              # Go 多模块工作区配置
├─ .gitignore
└─ README.md
```

---

## 🚀 快速开始

### 1️⃣ 克隆项目

```bash
git clone https://github.com/ser163/HAIR_CLONE.git
cd HAIR_CLONE
```

### 2️⃣ 初始化工作区

```bash
go work init ./agent/ipghost ./ip_pool_chick ./ip_mole
go work sync
```

### 3️⃣ 安装依赖（如首次 clone）

```bash
go mod tidy ./agent/ipghost
go mod tidy ./ip_pool_chick
go mod tidy ./ip_mole
```

---

## 🛠️ 各模块说明

### `ip_pool_chick`（代理 IP 池）

- 功能：维护 Redis 堆栈结构的 IP 池，最大容量 100 个，动态弹出推入。
- 特点：使用 Redis 管理数据池，后端数据源为 MySQL 表 `checklist`。
- 技术：Go + Redis + MySQL

### `ip_mole`（代理可用性检查器）

- 功能：每 30 分钟定时从数据库读取待检测 IP，校验端口和 HTTP 代理功能是否可用。
- 策略：
    - 检查失败数 > 3 的记录自动删除；
    - 更新 `fail`, `succ`, `check_num`, `last_check` 等字段；
- 技术：Go + MySQL + HTTP

### `agent/ipghost`（IP 收集器）

- 功能：从不同渠道抓取代理 IP 并写入数据库或消息队列（待扩展）。
- 特点：可部署多个不同策略的收集器（HTTP 抓取 / API / 采集器接入）。

---

## ⚙️ 配置说明

每个模块的 `config/config.yaml` 管理数据库、检测超时等参数：

```yaml
database:
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "password"
  dbname: "ip_pool"
  charset: "utf8mb4"

check:
  portTimeout: 3s
  httpTimeout: 5s
  targetURL: "http://example.com"
```

---

## 🔁 运行方式示例

运行检测模块：

```bash
cd ip_mole
go run cmd/main.go
```

运行池管理模块：

```bash
cd ip_pool_chick
go run cmd/main.go
```

---

```bash
go test ./...
```

---

## 📝 License

MIT License © 2025 by Harry Liu

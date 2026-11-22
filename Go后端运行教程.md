# Go 后端启动教程

启动后端只需要完成下面 5 步：

1. 安装 Go
2. 启动 MySQL
3. 创建后端配置文件
4. 创建数据库
5. 下载依赖并启动

项目目录：

```text
2027fe_homework_be/   后端
fe-vue-demo/          前端
```

后端地址：`http://127.0.0.1:8088`

接口前缀：`/api`

---

## 1. 安装 Go

打开 https://go.dev/dl/，下载适合当前系统的 Go 安装包。

- Windows：下载 `.msi`，双击安装并保持默认选项。
- macOS：下载 `.pkg`，或使用 `brew install go`。
- Linux：可使用发行版包管理器安装 Go。

安装后重新打开终端，运行：

```bash
go version
```

版本需要是 `1.25` 或更高，例如 `go1.25.0`。

如果找不到 `go` 命令，先关闭并重新打开终端，然后确认 Go 已加入系统 `PATH`。

---

## 2. 启动 MySQL

安装 MySQL 8，记住设置的 `root` 密码，并保持端口为 `3306`。

Windows 打开 PowerShell，检查并启动 MySQL 服务：

```powershell
Get-Service
Start-Service MySQL80
```

如果实际服务名不是 `MySQL80`，把 `MySQL80` 换成服务列表中 MySQL 对应的名称。

验证登录：

```bash
mysql -u root -p
```

输入密码后看到 `mysql>` 即表示正常。输入 `exit;`退出。

如果终端找不到 `mysql`，改用下面的完整路径登录：

```powershell
& "C:\Program Files\MySQL\MySQL Server 8.0\bin\mysql.exe" -u root -p
```

---

## 3. 创建后端配置

在项目根目录打开终端，进入后端目录：

```bash
cd 2027fe_homework_be
```

复制配置模板：

**Windows PowerShell：**

```powershell
Copy-Item .env.example .env
```

**macOS / Linux：**

```bash
cp .env.example .env
```

用记事本、VS Code 或其他编辑器打开 `.env`，至少修改数据库密码：

```env
APP_MYSQL_HOST=127.0.0.1
APP_MYSQL_PORT=3306
APP_MYSQL_NAME=be_exam
APP_MYSQL_USER=root
APP_MYSQL_PASS=你的MySQL密码

APP_ALLOW_ORIGINS=http://127.0.0.1:5173
```

不要把真实密码提交到 Git。`.env` 是本地配置文件，不需要提交。

---

## 4. 创建数据库

仍然在 `2027fe_homework_be/` 目录执行：

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS be_exam CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

输入 MySQL 密码。

确认数据库已创建：

```bash
mysql -u root -p -e "SHOW DATABASES LIKE 'be_exam';"
```

输出中看到 `be_exam` 即表示成功。

不需要手动导入建表 SQL。后端启动时会在开发模式下使用 GORM 自动创建用户、文章、评论等表。不要在 `.env` 中设置 `APP_PROD`。

---

## 5. 下载依赖并启动

在 `2027fe_homework_be/` 目录执行：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
go run .
```

第一次下载依赖可能需要几分钟。

启动成功后，终端会持续运行并输出日志，不要按 `Ctrl+C`。

停止后端：在运行后端的终端按 `Ctrl+C`。

---

## 验证后端

重新打开一个终端。

**浏览器直接打开：**

```text
http://127.0.0.1:8088/api/articles/list
```

**PowerShell：**

```powershell
Invoke-RestMethod "http://127.0.0.1:8088/api/articles/list"
```

看到 `code: 0` 或页面接口数据即表示后端启动成功。

---

## 启动前端

保持后端终端运行，另开终端：

```bash
cd fe-vue-demo
pnpm install
pnpm dev
```

浏览器打开：

```text
http://127.0.0.1:5173
```

前端和后端都在本机时无需修改前端端口。

---

## 常见报错

### `connection refused` 或连接不到 `127.0.0.1:3306`

MySQL 未启动：

```powershell
Get-Service
Start-Service MySQL80
```

### `Access denied for user 'root'`

`.env` 中的 `APP_MYSQL_PASS` 与 MySQL 的实际密码不一致。逐字修改后重启后端。

### `Unknown database 'be_exam'`

重新创建数据库：

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS be_exam CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### `Table doesn't exist`

检查：

- `.env` 中不能有 `APP_PROD`
- 当前终端必须位于 `2027fe_homework_be/` 目录
- 数据库确实是 `be_exam`

修改后重新执行 `go run .`。

### `bind: address already in use`

8088 端口已被占用：

```powershell
Get-NetTCPConnection -LocalPort 8088
```

记录 `OwningProcess`，再结束对应进程：

```powershell
Stop-Process -Id 进程号
```

### `go: command not found` 或提示找不到 `go`

重新安装或配置 Go；Windows 安装后关闭并重新打开 VS Code 和终端。

### 启动时找不到 `.env`

确认已进入后端目录，并确认 `.env` 已创建：

```bash
cd 2027fe_homework_be
Get-Item .env
go run .
```

最后一步中的 `Get-Item` 仅适用于 PowerShell。

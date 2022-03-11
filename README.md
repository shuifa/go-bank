# go-bank project for study golang

- [x] 数据库表设计design网站 [dbdiagram.io](https://dbdiagram.io/d/6225e04d61d06e6eadb5d9cc)

- [x] 安装 docker `brew install docker --cask`，学习docker 镜像、容器、数据卷的使用。拉取Postgres镜像并启动容器运行

- [x] golang数据迁移包 [golang-migrate](https://github.com/golang-migrate/migrate) ，使用 migrate 管理数据迁移和版本切换

- [x] 使用 Makefile 创建项目命令

- [x] 生成 CRUD 代码包 [sqlc](https://github.com/kyleconroy/sqlc)  sqlc 的 yml文件配置

- [x] golang 单元测试包 [testify](https://github.com/stretchr/testify) `main_test`，`go test -v ./...`

- [x] golang操作事务封装，协程测试数据库并发，学习数据库死锁、锁超时。使用乐观锁更改余额。保证相同的操作顺序避免死锁。

- [x] 事务隔离级别（读未提交、读已提交、可重复读、序列化）分别会遇到什么问题（脏读、幻读、不可重复读、无法序列化）

- [x] 使用 GitHub Actions 在push代码时运行代码测试（CI），ci工作流 `job --> steps --> actions`。定义工作流的service和dependence

- [x] 开发`restful`风格的 `webapi`，[gin-web](https://github.com/gin-gonic/gin) ，`postman`测试api

- [x] 从文件、环境变量、配置中心读取配置的包 [viper](https://github.com/spf13/viper) 

- [x] mockdb 的包 [mock](https://github.com/golang/mock)，编写测试用例覆盖所有情况
# rop2-api
rop2-web的配套api。

### 开发/运行指南
需要安装MySQL 8.0。

需配置的环境变量(所有配置环境变量都以ROP2_开头，具体见utils/config.go)：
| 后缀(以ROP2_开头)  | 描述             | 格式                                                            |
| ------------------ | ---------------- | --------------------------------------------------------------- |
| Addr               | 监听地址         | 0.0.0.0:8080                                                    |
| DSN                | 数据库连接字符串 | username:pwd@tcp(host:port)/rop2?charset=utf8mb4&parseTime=true |
| IdentityKey        | token签名私钥    | 随机生成的任意文本                                              |
| LoginCallbackRegex | 登录回调地址正则 | ^http://localhost:5173(/.*)?$                                   |

### 注意事项
- gin框架的`BindJSON`不能使用，需使用`ShouldBindJSON`
- 使用`ShouldBindJSON`时，字段名需首字母大写，实际请求不区分大小写
- 使用`ShouldBindQuery`时，字段名需首字母大写，实际请求区分大小写，用`form:"name"`指定实际参数名
- JSON序列化时会将对象按键名排序
- map\[x\]y是无序的
- 导入包不能有循环依赖，如utils不能再导入model handler等
- gorm的bug|feature比你想象的要多。请谨慎使用`Save`。

### 部署方式
确保cwd下恰包含docker-compose.yml

```sh
docker-compose-18 stop
docker-compose-18 up --build -d
```
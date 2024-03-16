# rop2-api
rop2-web的配套api。

### 开发/运行指南
需要安装MySQL 8.0。

需配置的环境变量(所有配置环境变量都以ROP2_开头，具体见utils/config.go)：
| 后缀        | 描述             | 格式                                                            |
| ----------- | ---------------- | --------------------------------------------------------------- |
| DSN         | 数据库连接字符串 | username:pwd@tcp(host:port)/rop2?charset=utf8mb4&parseTime=true |
| IdentityKey | token签名私钥    | 随机生成的任意文本                                              |

### 注意事项
不灵仍在学习Golang开发，他将一些注意事项/容易踩的坑写在这里。
- gin框架的`BindJSON`不能使用，需使用`ShouldBindJSON`
- 使用`ShouldBindJSON`时，字段名需首字母大写，实际请求不区分大小写
- 使用`ShouldBindQuery`时，字段名需首字母大写，实际请求区分大小写，用`form:"name"`指定实际参数名
- JSON序列化时会将对象按键名排序
- map\[x\]y是无序的
- 导入包不能有循环依赖，如utils不能再导入model handler等
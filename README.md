# rop2-api
rop2-web的配套api。

### 注意事项
- gin框架的`BindJSON`不能使用，需使用`ShouldBindJSON`
- 使用`ShouldBindJSON`时，字段名需首字母大写，实际请求不区分大小写
- 使用`ShouldBindQuery`时，字段名需首字母大写，实际请求区分大小写，用`form:"name"`指定实际参数名
- JSON序列化时会将对象按键名排序
- map\[x\]y是无序的
- 导入包不能有循环依赖，如utils不能再导入model handler等
- gorm的bug|feature比你想象的要多。请谨慎使用`Save`。

### 部署方式
确保cwd下恰包含compose.yml。请使用Docker Compose v2.20+

```sh
docker-compose-2 build
docker-compose-2 stop
docker-compose-2 up --no-build -d
```
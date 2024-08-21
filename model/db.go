package model

import (
	"fmt"
	"rop2-api/utils"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// TODO 测试组织ID，会给所有登录用户添加权限，正式环境应删除
var TestOrgId uint32 = 1

func Init() {
	var err error
	db, err = gorm.Open(mysql.Open(utils.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		TranslateError:                           true,
		DisableAutomaticPing:                     false, //开启定期ping，防止无操作断连
	})
	if err != nil {
		panic(err)
	}
	var testOrgIds []uint32
	db.Model(&Org{}).Where("name = ?", "测试组织").Limit(1).Pluck("id", &testOrgIds)
	if len(testOrgIds) > 0 {
		TestOrgId = testOrgIds[0]
	} else {
		fmt.Println("未找到测试组织(name = ?)")
	}
}

// 删除并重建数据库结构
func ResetDb() {
	db.Exec("DROP DATABASE IF EXISTS rop2;")
	db.Exec("CREATE DATABASE rop2;")
	db.Exec("USE rop2;")
	migrator := db.Migrator()
	migrator.AutoMigrate(&Org{}, &Depart{}, &Admin{}, &Template{}, &Stage{}, &Form{}, &Person{}, &Intent{}, &Result{}, &Interview{},
		&InterviewSchedule{})

	//建表完成，添加外键
	const restrict, cascade, setNull = "RESTRICT", "CASCADE", "SET NULL"
	fkBuilder := func(thisTable, thisCol, refTable, refCol, onDelete string) string {
		lines := []string{
			fmt.Sprintf("ALTER TABLE `%s`", thisTable),
			fmt.Sprintf("ADD CONSTRAINT fk_%s_%s", thisTable, thisCol),
			fmt.Sprintf("FOREIGN KEY (`%s`)", thisCol),
			fmt.Sprintf("REFERENCES `%s` (`%s`)", refTable, refCol),
			"ON UPDATE CASCADE", //固定外键追随更新，这个行为一般不需要改变
			fmt.Sprintf("ON DELETE %s;", onDelete),
		}
		return strings.Join(lines, "\n")
	}

	db.Exec(fkBuilder("orgs", "default_depart", "departs", "id", cascade))

	db.Exec(fkBuilder("departs", "owner", "orgs", "id", cascade)) //删除组织时自动删除所有部门

	db.Exec(fkBuilder("admins", "at", "orgs", "id", cascade)) //删除组织时自动删除所有管理员

	db.Exec(fkBuilder("templates", "owner", "orgs", "id", cascade)) //删除组织时自动删除所有模板

	db.Exec(fkBuilder("stages", "owner", "departs", "id", cascade))      //删除部门时自动删除相关的阶段设定
	db.Exec(fkBuilder("stages", "on_enter", "templates", "id", setNull)) //删除通知模板不删除阶段设定

	db.Exec(fkBuilder("forms", "owner", "orgs", "id", cascade)) //删除组织时自动删除所有表单

	db.Exec(fkBuilder("intents", "zju_id", "people", "zju_id", cascade))
	db.Exec(fkBuilder("intents", "form", "forms", "id", cascade))

	db.Exec(fkBuilder("results", "zju_id", "people", "zju_id", cascade))
	db.Exec(fkBuilder("results", "form", "forms", "id", cascade))

	db.Exec(fkBuilder("interviews", "form", "forms", "id", cascade))
	db.Exec(fkBuilder("interviews", "depart", "departs", "id", cascade))

	db.Exec(fkBuilder("interview_schedules", "zju_id", "people", "zju_id", cascade))
	db.Exec(fkBuilder("interview_schedules", "interview", "interviews", "id", cascade))

	//数据库初始化完成，但不添加任何测试数据

	TestOrgId, _ := InitNewOrg("测试组织", "_", "测试管理员")
	CreateDepart(TestOrgId, "部门1")
}

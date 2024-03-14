package model

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	if dsn, ok := os.LookupEnv("ROP2_DSN"); ok {
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("dsn not found"))
	}
}

// 删除并重建数据库结构
func ResetDb() {
	db.Exec("DROP DATABASE IF EXISTS rop2;")
	db.Exec("CREATE DATABASE rop2;")
	db.Exec("USE rop2;")
	migrator := db.Migrator()
	migrator.AutoMigrate(&Org{}, &Depart{}, &Form{}, &User{})

	//建表完成，添加外键
	fkBuilder := func(thisTable, thisCol, refTable, refCol, onDelete string) string {
		lines := []string{
			fmt.Sprintf("ALTER TABLE `%s`", thisTable),
			fmt.Sprintf("ADD CONSTRAINT fk_%s_%s", thisTable, thisCol),
			fmt.Sprintf("FOREIGN KEY (`%s`)", thisCol),
			fmt.Sprintf("REFERENCES `%s` (`%s`)", refTable, refCol),
			"ON UPDATE CASCADE",
			fmt.Sprintf("ON DELETE %s;", onDelete),
		}
		return strings.Join(lines, "\n")
	}
	const restrict, cascade, setNull = "RESTRICT", "CASCADE", "SET NULL"
	db.Exec(fkBuilder("departs", "parent", "orgs", "id", restrict))         //删除组织前需删除所有部门
	db.Exec(fkBuilder("orgs", "default_depart", "departs", "id", restrict)) //默认部门绑定后不能删除
	db.Exec(fkBuilder("forms", "owner", "orgs", "id", restrict))            //删除组织前需删除所有表单
	db.Exec(fkBuilder("users", "at", "orgs", "id", cascade))                //删除组织时自动删除所有管理

	testOrg := &Org{
		Name: "测试组织",
	}
	db.Select("Name").Create(&testOrg)

	testOrgDefaultDepart := &Depart{
		Name:   "默认部门",
		Parent: testOrg.Id,
	}
	db.Select("Name", "Parent").Create(testOrgDefaultDepart)

	testOrg.DefaultDepart = testOrgDefaultDepart.Id
	db.Save(&testOrg)

	//TODO 考虑是否删除测试数据
}

package model

import (
	"fmt"
	"rop2-api/utils"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(mysql.Open(utils.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		TranslateError:                           true,
	})
	if err != nil {
		panic(err)
	}
}

var (
	TestOrg              *Org
	TestOrgDefaultDepart *Depart
	TestUser             *User
	TestStage            *Stage
	TestForm             *Form
)

// 删除并重建数据库结构
func ResetDb() {
	db.Exec("DROP DATABASE IF EXISTS rop2;")
	db.Exec("CREATE DATABASE rop2;")
	db.Exec("USE rop2;")
	migrator := db.Migrator()
	migrator.AutoMigrate(&Org{}, &Depart{}, &User{}, &Stage{}, &Form{}, &Candidate{})

	//建表完成，添加外键
	const restrict, cascade, setNull = "RESTRICT", "CASCADE", "SET NULL"
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

	db.Exec(fkBuilder("departs", "parent", "orgs", "id", cascade))          //删除组织时自动删除所有部门
	db.Exec(fkBuilder("orgs", "default_depart", "departs", "id", restrict)) //默认部门不能删除（只能随组织一起删除）

	db.Exec(fkBuilder("users", "at", "orgs", "id", cascade)) //删除组织时自动删除所有管理员

	db.Exec(fkBuilder("stages", "owner", "orgs", "id", cascade))  //删除组织时自动删除所有阶段
	db.Exec(fkBuilder("stages", "next", "stages", "id", setNull)) //删除阶段时，清除next绑定

	db.Exec(fkBuilder("forms", "owner", "orgs", "id", cascade))    //删除组织时自动删除所有表单
	db.Exec(fkBuilder("forms", "enter", "stages", "id", restrict)) //删除阶段时不能有表单绑定在其上

	db.Exec(fkBuilder("candidates", "of", "departs", "id", cascade))   //删除组织时，移除所有候选人
	db.Exec(fkBuilder("candidates", "form", "forms", "id", cascade))   //删除表单时，移除所有候选人
	db.Exec(fkBuilder("candidates", "stage", "stages", "id", cascade)) //删除阶段时，移除所有候选人
	db.Exec(fkBuilder("candidates", "gone", "candidates", "id", setNull))

	TestOrg = &Org{
		Name: "测试组织",
	}
	db.Select("Name").Create(TestOrg)

	TestOrgDefaultDepart = &Depart{
		Name:   "默认部门",
		Parent: TestOrg.Id,
	}
	db.Select("Name", "Parent").Create(TestOrgDefaultDepart)

	TestOrg.DefaultDepart = TestOrgDefaultDepart.Id
	db.Save(TestOrg)

	TestUser = &User{
		ZjuId:    "__N/A__",
		Nickname: "测试用户",
		At:       TestOrg.Id,
		Perm: utils.Stringify(PermMap{
			(TestOrgDefaultDepart.Id): Maintainer,
		}),
	}
	db.Select("ZjuId", "Nickname", "At", "Perm").Create(TestUser)

	TestStage = &Stage{
		Name:  "填表",
		Owner: TestOrg.Id,
		Tasks: "[]",
	}
	db.Select("Name", "Owner", "Tasks").Create(TestStage)

	TestForm = &Form{
		Name:     "测试组织2024年春季纳新报名表",
		Entry:    1,
		Children: `[{"id":1}]`,
		Owner:    TestOrg.Id,
		Enter:    TestStage.Id,
	}
	db.Select("Name", "Entry", "Children", "Owner", "Enter").Create(TestForm)

	//TODO 考虑是否删除测试数据
}

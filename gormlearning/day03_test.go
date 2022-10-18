package gormlearning

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
	"time"
)

type User struct {
	ID           uint `gorm:"primaryKey;not null;autoIncrement"`
	Name         string
	Email        *string //默认值为null
	Sex          string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreateAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `grom:"autoUpdateTime"`
}
type Email struct {
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.Name == "sb" {
		fmt.Println("不文明命名！")
	}

	return nil
}
func TestDay03_1(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gormlearning?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	fmt.Println("连接数据库", err)
	err = db.AutoMigrate(&User{})
	fmt.Println("创建表", err)
	bir := time.Date(2000, 7, 2, 3, 4, 5, 6, time.UTC)
	memnum := sql.NullString{
		String: "12343545",
		Valid:  true,
	}
	user := &User{Name: "o1", Age: 2, Birthday: &bir, MemberNumber: memnum}
	result := db.Create(&user)
	fmt.Println("创建模式userID", user.ID, "result", result.Error, result.RowsAffected)
	user1 := &User{Name: "o2", Age: 2, Birthday: &bir, MemberNumber: memnum}
	result = db.Omit("Age").Create(&user1)
	fmt.Println("忽略模式：userid", user.ID, "result", result.Error, result.RowsAffected)
	user3 := &User{Name: "sb", Age: 2, Birthday: &bir, MemberNumber: memnum}
	db.Create(&user3)
	db.Model(&User{})
	user4 := &User{Name: "sb", Age: 2, Sex: "女", Birthday: &bir, MemberNumber: memnum}
	db.Model(&user4).Create(map[string]interface{}{
		"Name": "o5",
		"Sex": clause.Expr{SQL: "?", Vars: []interface{}{
			// 性别字符串转int
			func(string2 string) int {
				if string2 == "男" {
					return 1
				}
				return 0
			}("女"),
		}},
	})
	result = db.First(&user)
	fmt.Println("查询第一行", result.RowsAffected, result.Error, user.ID)
	rows, _ := result.Rows()
	for rows.Next() {
		u := &User{}
		err = db.ScanRows(rows, u)
		fmt.Println("查询结果", u)
	}

}
func TestDay03_2(t *testing.T) {
	//bir := time.Date(2000, 7, 2, 3, 4, 5, 6, time.UTC)
	memnum := sql.NullString{
		String: "12343545",
		Valid:  true,
	}

	user := &User{}
	dsn := "root:123456@tcp(127.0.0.1:3306)/gormlearning?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	result := db.First(&user)
	fmt.Println("查询第一行", result.RowsAffected, result.Error, user.ID)
	rows, _ := result.Rows()
	for rows.Next() {
		u := &User{}
		_ = db.ScanRows(rows, u)
		fmt.Println("查询结果", u)
	}
	user1 := &User{}
	db.Last(&user1)
	fmt.Println("查询最后一个", user1)
	user2 := &User{} //不管user里放的什么数，都不会影响查询结果
	db.Take(&user2)
	fmt.Println("随便查询一条", user2)
	user3 := &User{}
	var userArr []User
	rows, _ = db.Where("name=?", "o1").Find(&userArr).Rows()
	for rows.Next() {
		u := &User{}
		_ = db.ScanRows(rows, u)
		fmt.Println("查询结果2", u)
	}
	fmt.Println("查询用户1", user3)
	user4 := &User{}
	db.Where(&User{Name: "o1", Age: 2, MemberNumber: memnum}).Find(&user4)
	fmt.Println("查询用户2", user4)

	db.Where(map[string]interface{}{"id": []int{1, 2, 3, 4, 5}}).Find(&userArr)
	fmt.Println("查询数组", userArr)

	user5 := &User{}
	db.First(&user5, "id=? and name=?", 25, "o1")
	fmt.Println("查询结果user5", user5)

	var users1 []User
	var a int
	db.Model(&User{}).Select("sum(age)").Find(&a)
	fmt.Println("select", a)
	db.Order("age desc").Order("name").Find(&users1)
	fmt.Println("排序查询", users1)

	user6 := &User{}
	db.Raw("select id,name ,age from users where id=?", 5).Scan(user6)
	fmt.Println("user6", user6)
	var user7 []User

	db.Scopes(GrownUpAge).Find(&user7)
	fmt.Println("user7", user7)

	for i := 1; i <= 3; i++ {
		var user8 []User
		db.Scopes(PageHelper(i, 10)).Find(&user8)
		fmt.Printf("第%v页:%v", i, user8)
		fmt.Println()
	}

}
func GrownUpAge(db *gorm.DB) *gorm.DB {
	return db.Where("id >= ?", 5)
}
func PageHelper(pageNum int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func TestDay03_3(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gormlearning?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	user := &User{ID: 28, Name: "yj", Age: 5}
	//user1 := &User{ID: 29, Name: "haha", Age: 6}
	db.Model(&user).Select("Name", "Age").Updates(&user)
	//db.Save(user1)

}
func TestDay03_4(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gormlearning?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var users []User
	db.Raw("select id ,name,age from users where name=?", "o1").Scan(&users)
	fmt.Println("users", users)
	user := &User{}
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	fmt.Println("sql", stmt.SQL.String(), stmt.Vars)
	db.Transaction(func(tx *gorm.DB) error { //外层事务影响内层事务，相同级别不影响。
		user1 := &User{Name: "op", Age: 9}
		tx.Create(&user1)
		tx.Transaction(func(tx2 *gorm.DB) error {
			user2 := &User{Name: "xc", Age: 4}
			tx2.Create(&user2)
			return errors.New("rollback user2")
		})
		return errors.New("xss")
	})

}

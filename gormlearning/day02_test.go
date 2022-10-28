package gormlearning

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

type User1 struct {
	ID       int
	Name     string
	Age      int
	Birthday time.Time
}
type Pet struct {
	Name string
	Age  int
	ID   int
}
type Company struct {
	ID       uint `gorms:"primarykey""`
	Name     string
	Worker   []Worker `gorm:"foreignKey:ID;references:WorkerId"`
	WorkerId uint     `json:"projectId" gorm:"uique_index:fk_companies_worker"`
}
type Worker struct {
	ID         uint `gorms:"primarykey"`
	WorkerName string
	ComparyID  uint
}
type Location struct {
	ID   uint `gorms:"primarykey"`
	X, Y int
}
type Student struct {
	Name       string
	Location   Location `gorm:"foreignKey:LocationId"`
	LocationId int
}
type Stu struct {
	Name       string
	Locations  []Loca `gorm:"foreignKey:LocationId"`
	LocationId int
}
type Loca struct {
	ID   uint `gorms:"primarykey"`
	X, Y int
}

func TestDay02_1(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/users?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	fmt.Println("db", db, "err", err)
	bir := time.Date(2000, 7, 25, 0, 0, 0, 0, time.UTC)
	fmt.Println("utc", bir.UTC())
	//user := User{Name: "xm", Age: 23, Birthday: bir}
	//result := db.Create(&user)
	//fmt.Println("userid", user.ID, "result", result.Error, result.RowsAffected)
	//user1 := User{Name: "xh", Age: 22, Birthday: bir}
	//result := db.Select("Name", "birthday").Create(&user1)
	//fmt.Println("xbs result", result.Error, result.RowsAffected)
	//	db.Omit("Age").Create(&user)
	//	var users = []User{{Name: "jinzhu1", Age: 1, Birthday: bir}, {Name: "jinzhu2", Age: 2, Birthday: bir}, {Name: "jinzhu3", Age: 3, Birthday: bir}}
	//	db.Create(&users)
	//	for _, user = range users {
	//		fmt.Println("userid", user.ID)
	//	}
	//	user2 := User{Name: "jinzhu4", Age: 4, Birthday: bir, Pet: Pet{Name: "coco", Age: 1}}
	//	db.Create(&user2)
	//fmt.Println("user2", user2.ID)

	db.AutoMigrate(&Student{}, &Location{})

	result := db.Create(&Student{
		Name:     "jinzhu1",
		Location: Location{1, 100, 200},
	})
	fmt.Println("res", result.Error, result.RowsAffected)
	/*db.Model(Student{}).Create(map[string]interface{}{
		"Name":     "jinzhu",
		"Location": clause.Expr{SQL: "ST_PointFromText(?)", Vars: []interface{}{"POINT(100 100)"}},
	})*/
	/*db.Model(&User{}).Create(map[string]interface{}{
		"Name": "jinzhu5", "Age": 18, "birthday": bir,
	})*/

}

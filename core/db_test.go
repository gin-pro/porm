package core

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestScan(t *testing.T) {
	db, err := Open("mysql", "root:jimbir8520@tcp(203.189.212.101:7001)/porm?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(fmt.Sprintf("open err : %v", err))
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(fmt.Sprintf("Begin err : %v", err))
		return
	}
	rs, err := tx.Query(`select * from user`)
	if err != nil {
		fmt.Println(fmt.Sprintf("Query err : %v", err))
		return
	}

	fmt.Println(fmt.Sprintf("Res : %v", rs))
	for rs.Next() {
		var ID int
		var Name string
		var Sex string
		err = rs.Scan(&ID, &Name, &Sex)
		if err != nil {
			fmt.Println(fmt.Sprintf("Scan err : %v", err))
			return
		}
		fmt.Println(fmt.Sprintf("User : %v,%v,%v", ID, Name, Sex))
	}
}

func TestScanStructByIndex(t *testing.T) {

	type User struct {
		ID   int    `porm:"id"`
		Name string `porm:"name"`
		Sex  string `porm:"sex"`
	}

	db, err := Open("mysql", "root:jimbir8520@tcp(203.189.212.101:7001)/porm?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(fmt.Sprintf("open err : %v", err))
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(fmt.Sprintf("Begin err : %v", err))
		return
	}
	rs, err := tx.Query(`select * from user`)
	if err != nil {
		fmt.Println(fmt.Sprintf("Query err : %v", err))
		return
	}

	fmt.Println(fmt.Sprintf("Res : %v", rs))
	for rs.Next() {
		var user User
		err = rs.ScanStructByIndex(&user)
		if err != nil {
			fmt.Println(fmt.Sprintf("Scan err : %v", err))
			return
		}
		fmt.Println(fmt.Sprintf("User : %v", user))
	}
}

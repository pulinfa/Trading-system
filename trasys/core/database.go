package core

/*
*	数据库模块
 */

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func InitDB() (err error) {
	dsn := "root:88888888@tcp(127.0.0.1:3306)/my_db"
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("open database err, ", err)
		return err
	}

	err = Db.Ping()
	if err != nil {
		fmt.Println("ping database err, ", err)
		return err
	}

	return nil
}

//对于股票信息，只有一个功能就是查找股票的价格，所以这个相对好弄一点

// 查询所有的股票的信息
func QueryAllStock() {
	s := "select * from stockinfo;"
	r, err := Db.Query(s)

	var et string
	var sc string
	var lp float64
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		for r.Next() {
			r.Scan(&et, &sc, &lp)

			st := NewStock(et, sc, lp)

			GlobalStocks.AddStock(st)
		}
	}
}

//+++++++++++++++++++++++++++++++++++++++++++++++++++++
//对于持股信息

// 查询所有的持仓信息
func QueryAllReocrd() {
	s := "select * from holdinginfo;"
	r, err := Db.Query(s)

	var et string
	var sc string
	var lp float64
	var cid string
	var ha int32
	var mv float64

	if err != nil {
		fmt.Println("err: ", err)
	} else {
		for r.Next() {
			r.Scan(&cid, &et, &sc, &lp, &ha, &mv)

			st := NewStock(et, sc, lp)

			hi := NewHoldingInfo(cid, *st, ha)

			GlobalHoldingRecord.AddRecord(hi)
		}
	}
}

// 插入一条新的数据
func InsertRecord(cid string, et string, sc string, lp float64, ha int32, mv float64) {
	s := "insert into holdinginfo (clientid, exchangetype, stockcode, lastprice, holdamount, marketvalue) values(?, ?, ?, ?, ?, ?);"
	_, err := Db.Exec(s, cid, et, sc, lp, ha, mv)
	if err != nil {
		fmt.Println("insert err, ", err)
		return
	}

	fmt.Println("--------------------> insert succ.")
}

// 更新某条数据
func UpdateRecord(cid string, et string, sc string, ha int32, mv float64) {
	s := "update holdinginfo set holdamount = ?, marketvalue = ? where clientid = ? and exchangetype = ? and stockcode = ?;"
	_, err := Db.Exec(s, ha, mv, cid, et, sc)

	if err != nil {
		fmt.Println("update err, ", err)
		return
	}

	fmt.Println("--------------> update succ.")
}

// 删除数据
func DeleteRecord(cid string, et string, sc string) {
	s := "delete from holdinginfo where clientid = ? and exchangetype = ? and stockcode = ?;"
	_, err := Db.Exec(s, cid, et, sc)

	if err != nil {
		fmt.Println("delete err, ", err)
		return
	}

	fmt.Println("-------------------> delete succ.")
}

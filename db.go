package main

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	db, _ = sqlx.Connect("mysql", "root:123456@/jiyi")
	db.MapperFunc(strings.TrimSpace)
}

type Card struct {
	CardId   int
	Front    string
	Back     string
	Category int
}

func insertCard(front, back string) {
	sql := "INSERT INTO card(front, back, category) VALUES(? ,? ,1)"
	db.Exec(sql, front, back)
}

func selectCard() (res []Card) {
	sql := "SELECT * FROM card"
	db.Select(&res, sql)
	return
}

func getCard(cardId string) (res Card) {
	sql := "SELECT * FROM card WHERE cardId=?"
	db.Get(&res, sql, cardId)
	return
}

package main

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	db, _ = sqlx.Connect("mysql", "root:123456@/jiyi")
	db.MapperFunc(strcase.ToSnake)
}

// 卡片
type Card struct {
	CardId int
	KindId int
	Front  string
	Back   string
	Helper string
	Pinyin string
}

func (c *Card) get(cardId int) *Card {
	sql := `SELECT * FROM card WHERE card_id=?`
	db.Get(c, sql, cardId)
	return c
}

func (c *Card) insert() {
	sql := `INSERT INTO card(kind_id, front, back, helper, pinyin) VALUES(? ,? ,?, ?, ?)`
	db.Exec(sql, c.KindId, c.Front, c.Back, c.Helper, c.Pinyin)
}

func (c *Card) update() {
	sql := `UPDATE card SET front=?, back=?, helper=?, pinyin=? WHERE card_id=?`
	db.Exec(sql, c.Front, c.Back, c.Helper, c.Pinyin, c.CardId)
}

func (c *Card) delete(cardId int) {
	sql := `DELETE FROM card WHERE card_id=?`
	db.Exec(sql, cardId)
}

// 删除卡片（todo：事物）
//
// 1.删除卡片
// 2.删除卡片与卡组的关联
func (c *Card) deleteTx(cardId int) {
	c.delete(cardId)
	cardDeck.deleteByCardId(cardId)
}

type Cards []Card

// todo：分页
func (c *Cards) list(kindId int) *Cards {
	sql := `SELECT * FROM card WHERE kind_id=? LIMIT 100`
	db.Select(c, sql, kindId)
	return c
}

func (c *Cards) search(kindId int, query string) *Cards {
	fronts := splitSpace(query)
	sql := `SELECT * FROM card WHERE kind_id=? AND front IN(?)`
	sql, args, _ := sqlx.In(sql, kindId, fronts)
	db.Select(c, sql, args...)
	return c
}

func (c *Cards) selectCardIds(kindId int, front string) *Cards {
	frontArray := splitSpace(front)
	sql := `SELECT * FROM card WHERE kind_id=? AND front IN(?) ORDER BY FIELD(front, ?)`
	sql, args, _ := sqlx.In(sql, kindId, frontArray, frontArray)
	db.Select(c, sql, args...)
	return c
}

func (c *Cards) selectFronts(deckId int) (fronts string) {
	sql := `SELECT card.* FROM card, deck, card_deck 
			WHERE deck.deck_id=? 
			AND card.card_id=card_deck.card_id 
			AND deck.deck_id=card_deck.deck_id
			ORDER BY card_deck_id`
	db.Select(c, sql, deckId)

	for _, v := range *c {
		fronts += v.Front + "\n"
	}

	return
}

// 卡组
type Deck struct {
	DeckId   int
	KindId   int
	DeckName string
}

func (d *Deck) get(deckId int) *Deck {
	sql := `SELECT * FROM deck WHERE deck_id=?`
	db.Get(d, sql, deckId)
	return d
}

func (d *Deck) insert() int {
	sql := `INSERT INTO deck(deck_name, kind_id) VALUES(?, ?)`
	res, _ := db.Exec(sql, d.DeckName, d.KindId)
	lastId, _ := res.LastInsertId()
	return int(lastId)
}

func (d *Deck) update() {
	sql := `UPDATE deck SET deck_name=? WHERE deck_id=?`
	db.Exec(sql, d.DeckName, d.DeckId)
}

func (d *Deck) delete(deckId int) {
	sql := `DELETE FROM deck WHERE deck_id=?`
	db.Exec(sql, deckId)
}

// 添加卡组（todo：事物）
//
// 1.添加卡组
// 2.获取卡片ID
// 3.添加卡片与卡组的关联
func (d *Deck) insertTx(fronts string) {
	deckId := d.insert()
	cardArray := cards.selectCardIds(d.KindId, fronts)
	cardDeck.insert(cardArray, deckId)
}

// 更新卡组（todo：事物）
//
// 1.更新卡组
// 2.删除旧的关联
// 3.获取卡片ID
// 4.添加卡片与卡组的关联
func (d *Deck) updateTx(fronts string) {
	d.update()
	cardDeck.deleteByDeckId(d.DeckId)
	c := cards.selectCardIds(d.KindId, fronts)
	cardDeck.insert(c, d.DeckId)
}

// 删除卡组（todo：事物）
//
// 1.删除卡组
// 2.删除卡片与卡组的关联
func (deck *Deck) deleteTx(deckId int) {
	deck.delete(deckId)
	cardDeck.deleteByDeckId(deckId)
}

type Decks []Deck

func (d *Decks) list(kindId int) *Decks {
	sql := `SELECT * FROM deck WHERE kind_id=? LIMIT 100`
	db.Select(d, sql, kindId)
	return d
}

// 卡片卡组关联
type CardDeck struct {
	CardDeckId int
	CardId     int
	DeckId     int
}

// 卡片与卡组的关联操作
func (cd *CardDeck) insert(c *Cards, deckId int) {
	cardDeckArray := []map[string]interface{}{}
	for _, v := range *c {
		cardId := v.CardId
		row := map[string]interface{}{"cardId": cardId, "deckId": deckId}
		cardDeckArray = append(cardDeckArray, row)
	}
	sql := `INSERT INTO card_deck (card_id, deck_id) VALUES (:cardId, :deckId)`
	db.NamedExec(sql, cardDeckArray)
}

func (cardDeck *CardDeck) deleteByCardId(cardId int) {
	sql := `DELETE FROM card_deck WHERE card_id=?`
	db.Exec(sql, cardId)
}

func (cardDeck *CardDeck) deleteByDeckId(deckId int) {
	sql := `DELETE FROM card_deck WHERE deck_id=?`
	db.Exec(sql, deckId)
}

// 类型
type Kind struct {
	KindId   int
	KindName string
}

func (kind *Kind) get(kindId int) *Kind {
	kind.KindId = kindId
	switch kindId {
	case 1:
		kind.KindName = "普通"
	case 2:
		kind.KindName = "汉字"
	case 3:
		kind.KindName = "词语"
	case 4:
		kind.KindName = "古诗文"
	}
	return kind
}

// 公用函数
func splitSpace(s string) (res []string) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	res = strings.Fields(s)
	return
}

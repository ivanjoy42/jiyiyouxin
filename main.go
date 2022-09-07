package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("tpl/*")

	r.GET("/", index)
	card(r.Group("card"))
	deck(r.Group("deck"))

	r.Run(":8080")
}

func index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

// 卡片操作
//
// list列表页面；
// modify修改页面，update更新记录；
// create新建页面，insert插入记录；
// delete删除记录。
//
// todo：搜索卡片、分组浏览
func card(r *gin.RouterGroup) {
	r.GET("list", func(c *gin.Context) {
		card := selectCard()
		c.HTML(200, "cardList.html", gin.H{"card": card})
	})

	r.GET("modify", func(c *gin.Context) {
		cardId := c.Query("cardId")
		card := getCard(cardId)
		c.HTML(200, "cardModify.html", gin.H{"card": card})
	})

	r.GET("update", func(c *gin.Context) {
		cardId := c.Query("cardId")
		front := c.Query("front")
		back := c.Query("back")
		updateCard(cardId, front, back)
	})

	r.GET("create", func(c *gin.Context) {
		c.HTML(200, "cardCreate.html", gin.H{})
	})

	r.GET("insert", func(c *gin.Context) {
		front := c.Query("front")
		back := c.Query("back")
		insertCard(front, back)
	})

	r.GET("delete", func(c *gin.Context) {
		cardId := c.Query("cardId")
		deleteCard(cardId)
	})
}

// 卡组操作
//
// list列表页面；
// modify修改页面，update更新记录；
// create新建页面，insert插入记录；
// delete删除记录。
//
// todo：改为post
func deck(r *gin.RouterGroup) {
	r.GET("list", func(c *gin.Context) {
		deck := selectDeck()
		c.HTML(200, "deckList.html", gin.H{"deck": deck})
	})

	r.GET("modify", func(c *gin.Context) {
		deckId := c.Query("deckId")
		deck := getDeck(deckId)
		card := selectCardByDeckId(deckId)
		c.HTML(200, "deckModify.html", gin.H{"deck": deck, "card": card})
	})

	r.GET("update", func(c *gin.Context) {
		deckId := c.Query("deckId")
		deckName := c.Query("deckName")
		kind := c.Query("kind")
		cards := c.Query("cards")
		updateDeck(deckId, deckName)
		updateCardDeck(deckId, kind, cards)
	})

	r.GET("create", func(c *gin.Context) {
		c.HTML(200, "deckCreate.html", gin.H{})
	})

	r.GET("insert", func(c *gin.Context) {
		deckName := c.Query("deckName")
		insertDeck(deckName)
	})

	r.GET("delete", func(c *gin.Context) {
		deckId := c.Query("deckId")
		deleteDeck(deckId)
	})
}

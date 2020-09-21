package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

//sqlite使用

//モデルの定義。gorm.Modelで標準のモデルを呼び出し
type Todo struct {
	gorm.Model
	Text   string
	Status string
}

// テーブル名：todos
// カラム: text status

//DB初期化
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInit）")
	}

	db.AutoMigrate(&Todo{})
	//ここでは、ファイルが無ければ生成を行い、
	//すでにファイルがありマイグレートも行われていれば何も行いません。
	defer db.Close()
}

//DB追加
func dbInsert(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbInsert)")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

//DB更新
func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbUpdate)")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

//DB削除
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！（dbDelete)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

//DB全取得
func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetAll())")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

//DB一つ取得
func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開けず！(dbGetOne())")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

func main() {
	// defer db.SqlDB.Close()

	dbInit()

	router := initRouter()

	router.Run(":8008")
}
func initRouter() *gin.Engine {
	router := gin.Default()
	// router := gin.New()

	router.LoadHTMLGlob("*")

	router.GET("/", index)

	router.GET("/db", func(ctx *gin.Context) {
		todos := dbGetAll()
		fmt.Println(todos)
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
	})
	//Create
	router.POST("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbInsert(text, status)
		ctx.Redirect(302, "/")
	})

	//Update
	router.POST("/updata/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		fmt.Println(n)
		id, err := strconv.Atoi(n)
		fmt.Println(id)
		if err != nil {
			panic("ERROR")
		}
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbUpdate(id, text, status)
		ctx.Redirect(302, "/")
	})

	//Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		ctx.Redirect(302, "/")

	})

	router.POST("/", indexPost)

	router.POST("/test", indexLoop)

	return router
}

func helloWorld(context *gin.Context) {
	context.JSON(http.StatusOK, "Hello World")
}

func index(context *gin.Context) {
	// context.HTML(http.StatusOK, "index.html", "")
	context.HTML(http.StatusOK, "index.html", gin.H{
		"username": "",
		"text":     "",
	})
}
func indexPost(e04 *gin.Context) {
	username := e04.PostForm("username")
	fmt.Println(username)
	e04.HTML(http.StatusOK, "index.html", gin.H{
		"username": username,
	})
}

func indexLoop(context *gin.Context) {
	text := context.PostForm("text")
	fmt.Println(text)
	for i := 0; i < 100000; i++ {
		context.JSON(http.StatusOK, text)

	}
}

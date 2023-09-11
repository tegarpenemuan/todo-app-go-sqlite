package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" //mode local/development
)

var db *gorm.DB

type Todo struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	Title  string `json:"title"`
	Completed bool `json:"completed"`
}

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "./db/todo.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Todo{})

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		var todos []Todo
		db.Find(&todos)
		c.HTML(200, "todo.html", gin.H{
			"todos": todos,
		})
		fmt.Println(todos)
	})

	router.POST("/add", func(c *gin.Context) {
		title := c.PostForm("title") // Get the "title" field from the form
	
		if title != "" {
			todo := Todo{
				Title: title,
				Completed: false, // Assuming a new task is not completed initially
			}
			db.Create(&todo)
		}
	
		// Redirect to the home page or return a response as needed
		// You can also choose to return a JSON response with the newly created TODO task if you prefer.
		c.Redirect(http.StatusSeeOther, "/") // Redirect to the home page
	})

	router.POST("/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		var todo Todo
		db.First(&todo, id)
		todo.Completed = !todo.Completed
		db.Save(&todo)
		// c.JSON(200, todo)
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
		var todo Todo
		db.First(&todo, id)
		todo.Completed = !todo.Completed
		db.Delete(&todo)
		// c.JSON(200, todo)
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8080")
}


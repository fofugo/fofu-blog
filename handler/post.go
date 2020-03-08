package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v10"
)

type post struct {
	Title   string
	Content string
}

func PostHandler(c echo.Context) error {
	type queryParam struct {
		Id int `validate:"required" query:"id"`
	}
	input := queryParam{}
	if err := c.Bind(&input); err != nil {
		fmt.Println("error1")
	}
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		fmt.Println("error2")
	}

	db := c.Get("db").(*sql.DB)
	if err := db.Ping(); err != nil {
		panic(err)
	}

	rows, _ := db.Query("SELECT title,content FROM board WHERE id=?", input.Id)
	post := post{}
	for rows.Next() {
		rows.Scan(&post.Title, &post.Content)
	}
	return c.Render(http.StatusOK, "post", post)
}

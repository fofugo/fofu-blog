package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v10"
)

type Post struct {
	Title   string
	Content string
}

func PostHandler(c echo.Context) error {
	type queryParam struct {
		Id int `validate:"required" query:"id"`
	}
	input := queryParam{}
	if err := c.Bind(&input); err != nil {
		panic(err)
	}
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		panic(err)
	}

	db := c.Get("db").(*sql.DB)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	post := Post{}
	err := db.QueryRow("SELECT title,content FROM board WHERE id=?", input.Id).Scan(&post.Title, &post.Content)
	if err != nil {
		panic(err)
	}
	return c.Render(http.StatusOK, "post", post)
}

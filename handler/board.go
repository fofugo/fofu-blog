package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v10"
)

type board struct {
	Id       int
	Category string
	Title    string
}

func HomeHandler(c echo.Context) error {
	type queryParam struct {
		Page int `validate:"max=100,min=0" query:"page"`
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

	rows, _ := db.Query("SELECT id,category,title FROM board LIMIT 10 OFFSET ?", 10*input.Page)
	var boards []board
	for rows.Next() {
		board := board{}
		rows.Scan(&board.Id, &board.Category, &board.Title)
		boards = append(boards, board)
	}
	return c.Render(http.StatusOK, "Home", boards)
}

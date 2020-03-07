package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type Board struct {
	Id       int
	Category string
	Title    string
}

func BoardHandler(c echo.Context) error {
	page := c.QueryParam("page")
	fmt.Println(page)

	db := c.Get("db").(*sql.DB)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	rows, _ := db.Query("SELECT id,category,title FROM board")
	var boards []Board
	for rows.Next() {
		board := Board{}
		rows.Scan(&board.Id, &board.Category, &board.Title)
		boards = append(boards, board)
	}
	fmt.Println(boards)
	return c.Render(http.StatusOK, "board", boards)
}

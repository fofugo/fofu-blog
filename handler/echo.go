package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type echoTemplate struct {
	EchoBoardNos []echoBoardNo
	EchoBoards   []echoBoard
	Time         string
	Context      string
	CurrentPage  int
	Pages        []int
	No           int
}
type echoBoardNo struct {
	No      int
	Time    string
	Context string
}
type echoBoard struct {
	Title   string
	Content string
	Section string
}

func EchoHandler(c echo.Context) error {
	var echoTemplate echoTemplate
	type request struct {
		Page int `query:"page"`
		No   int `query:"no"`
	}
	input := request{}
	if err := c.Bind(&input); err != nil {
		panic(err)
	}
	if input.Page == 0 {
		input.Page = 1
	}
	echoTemplate.CurrentPage = input.Page
	var pages []int
	if input.Page == 1 {
		pages = []int{
			1,
			2,
			3,
		}
	} else {
		pages = []int{
			input.Page - 1,
			input.Page,
			input.Page + 1,
		}
	}
	echoTemplate.Pages = pages
	db := c.Get("db").(*sql.DB)
	if err := db.Ping(); err != nil {
		panic(err)
	}

	rows, _ := db.Query("SELECT no,time,context FROM  echo_board_no LIMIT 10 OFFSET ?", 5*(input.Page-1))
	for rows.Next() {
		var sql_time time.Time
		echoBoardNo := echoBoardNo{}
		rows.Scan(&echoBoardNo.No, &sql_time, &echoBoardNo.Context)
		echoBoardNo.Time = sql_time.Format("2006-01-02")
		echoTemplate.EchoBoardNos = append(echoTemplate.EchoBoardNos, echoBoardNo)
	}
	if input.No == 0 {
		input.No = echoTemplate.EchoBoardNos[0].No
	}
	var sql_time time.Time
	var context string
	if err := db.QueryRow("SELECT time,context FROM echo_board_no WHERE no = ?", input.No).Scan(&sql_time, &context); err != nil {
		panic(err)
	}
	echoTemplate.Time = sql_time.Format("2006-01-02")
	echoTemplate.Context = context

	rows, _ = db.Query("SELECT title,content FROM echo_board WHERE no = ?", input.No)
	count := 1
	for rows.Next() {
		echoBoard := echoBoard{}
		rows.Scan(&echoBoard.Title, &echoBoard.Content)
		echoBoard.Section = fmt.Sprintf("section%d", count)
		echoTemplate.EchoBoards = append(echoTemplate.EchoBoards, echoBoard)
		count++
	}
	echoTemplate.No = input.No

	return c.Render(http.StatusOK, "echo", echoTemplate)
}

package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type benchmarkTemplate struct {
	BenchmarkNos []benchmarkNo
	Benchmarks   []benchmark
	Time         string
	Context      string
	CurrentPage  int
	Pages        []int
	No           int
}
type benchmarkNo struct {
	No      int
	Time    string
	Context string
}
type benchmark struct {
	Title   string
	Content string
	Section string
}

func BenchmarkHandler(c echo.Context) error {
	var benchmarkTemplate benchmarkTemplate
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
	benchmarkTemplate.CurrentPage = input.Page
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
	benchmarkTemplate.Pages = pages
	db := c.Get("db").(*sql.DB)
	if err := db.Ping(); err != nil {
		panic(err)
	}

	rows, _ := db.Query("SELECT no,time,context FROM  benchmark_no LIMIT 10 OFFSET ?", 5*(input.Page-1))
	for rows.Next() {
		var sql_time time.Time
		benchmarkNo := benchmarkNo{}
		rows.Scan(&benchmarkNo.No, &sql_time, &benchmarkNo.Context)
		benchmarkNo.Time = sql_time.Format("2006-01-02")
		benchmarkTemplate.BenchmarkNos = append(benchmarkTemplate.BenchmarkNos, benchmarkNo)
	}
	if benchmarkTemplate.BenchmarkNos == nil {
		panic(errors.New("Error!"))
	}
	if input.No == 0 {
		input.No = benchmarkTemplate.BenchmarkNos[0].No
	}
	var sql_time time.Time
	var context string
	if err := db.QueryRow("SELECT time,context FROM benchmark_no WHERE no = ?", input.No).Scan(&sql_time, &context); err != nil {
		panic(err)
	}
	benchmarkTemplate.Time = sql_time.Format("2006-01-02")
	benchmarkTemplate.Context = context

	rows, _ = db.Query("SELECT title,content FROM benchmark WHERE no = ?", input.No)
	count := 1
	for rows.Next() {
		benchmark := benchmark{}
		rows.Scan(&benchmark.Title, &benchmark.Content)
		benchmark.Section = fmt.Sprintf("section%d", count)
		benchmarkTemplate.Benchmarks = append(benchmarkTemplate.Benchmarks, benchmark)
		count++
	}
	benchmarkTemplate.No = input.No
	return c.Render(http.StatusOK, "benchmark", benchmarkTemplate)
}

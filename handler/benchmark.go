package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v10"
)

type benchmarkTemplate struct {
	Context    string
	Benchmarks []Benchmark
}
type Benchmark struct {
	Title   string
	Content string
	Section string
}

func BenchmarkHandler(c echo.Context) error {
	type queryParam struct {
		Id int `validate:"required" query:"id"`
	}
	var benchmarkTemplate benchmarkTemplate
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
	err := db.QueryRow("SELECT context FROM benchmark_id WHERE id=?", input.Id).Scan(&benchmarkTemplate.Context)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT title,content FROM benchmark WHERE id=?", input.Id)
	if err != nil {
		panic(err)
	}
	count := 1
	for rows.Next() {
		benchmark := Benchmark{}
		rows.Scan(&benchmark.Title, &benchmark.Content)
		benchmark.Section = fmt.Sprintf("section%d", count)
		benchmarkTemplate.Benchmarks = append(benchmarkTemplate.Benchmarks, benchmark)
		count++
	}
	fmt.Println(benchmarkTemplate)
	return c.Render(http.StatusOK, "benchmark", benchmarkTemplate)
}

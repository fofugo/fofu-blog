package main

import (
	"fofu-blog/config"
	"fofu-blog/handler"
	"fofu-blog/middleware"

	"io"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// ---------------------------------------------------------------------

func main() {
	e := echo.New()

	dbConfig := config.DbConfig{
		Dialect:  "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "dongjulee",
		Password: "djfrnf081@",
		Name:     "fofu",
		Charset:  "utf8",
	}
	db := config.DB{}
	if err := db.Initialize(dbConfig); err != nil {
		return
	}
	e.Use(middleware.DbMiddleware(db))

	template := &Template{
		templates: template.Must(template.ParseGlob("gohtml/*.gohtml")),
	}
	e.Renderer = template
	// 첫 화면
	e.Static("/", "assets")
	e.GET("/", handler.HomeHandler)
	e.GET("/echo", handler.EchoHandler)
	e.GET("/benchmark", handler.BenchmarkHandler)
	e.Logger.Fatal(e.Start(":80"))
}

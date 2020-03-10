package main

import (
	"fofu-blog/handler"
	"fofu-blog/middleware"
	"fofu-blog/config"
	"os"
	"time"

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
	t := time.Now()
	startTime := t.Format("2006_01_02__15_04")
	f, err := os.OpenFile("log/log_"+startTime+".log", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	if err != nil {
		return
	}
	e.Logger.SetOutput(f)

	dbConfig := config.DbConfig{
		Dialect:  "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "",
		Password: "",
		Name:     "",
		Charset:  "",
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
	e.GET("/board", handler.BoardHandler)
	e.GET("/post", handler.PostHandler)
	e.GET("/benchmark/:id", handler.BenchmarkHandler)
	e.Logger.Fatal(e.Start(":80"))
}

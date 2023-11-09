package handlers

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gromples/mysqldb"
)

const STATIC_URL string = "/assets/"
const STATIC_ROOT string = "assets/"

type imageStruct struct {
	App_No    int
	Id        int
	Position  int
	ShortDesc string
	LongDesc  string
	Url       string
}

var tmpl = template.Must(template.ParseGlob("form/*.html"))

func ShowTemplateWithGallery(TemplateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Show Template With Gallery->", TemplateName)
		App_No := r.FormValue("App_No")
		if App_No == "" { //Default
			App_No = "1"
		}
		i, err := strconv.Atoi(App_No)
		if err != nil {
			i = 1
		}

		db := db.DBConn()
		defer db.Close()

		log.Println("Select from DB")
		selDB, err := db.Query("SELECT * FROM images where App_No = ? ORDER BY Position ASC", i)
		if err != nil {
			panic(err.Error())
		}
		img := imageStruct{}
		images := []imageStruct{}
		for selDB.Next() {
			var App_No int
			var Id int
			var Position int
			var ShortDesc string
			var LongDesc string
			var Url string
			err = selDB.Scan(&App_No, &Id, &Position, &ShortDesc, &LongDesc, &Url)
			if err != nil {
				panic(err.Error())
			}
			img.App_No = App_No
			img.Id = Id
			img.Position = Position
			img.ShortDesc = ShortDesc
			img.LongDesc = LongDesc
			img.Url = Url
			images = append(images, img)
		}
		log.Println("Before Template")
		tmpl.ExecuteTemplate(w, TemplateName, images)
		log.Println("After Template")
	}
}

func ShowTemplate(TemplateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Show Template ->", TemplateName)
		tmpl.ExecuteTemplate(w, TemplateName, nil)
	}
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

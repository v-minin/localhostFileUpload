package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var tpl *template.Template

func init() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	tpl = template.Must(template.ParseGlob("./templates/*"))
}

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/upload", upload)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":80", nil))
}

func foo(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Printf("template err [foo]: %s\n", err)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("parse form err %s\n", err)
		}

		mf, fh, err := r.FormFile("file")
		if err != nil {
			log.Printf("file err %s\n", err)
		}
		defer mf.Close()

		wd, _ := os.Getwd()
		path := filepath.Join(wd, "files", fh.Filename)
		nf, err := os.Create(path)
		if err != nil {
			log.Printf("file creation err %s\n", err)
		}
		defer nf.Close()

		mf.Seek(0, 0)
		io.Copy(nf, mf)

		err = tpl.ExecuteTemplate(w, "done.gohtml", fh.Filename)
		if err != nil {
			log.Printf("template err [done]: %s\n", err)
		}
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

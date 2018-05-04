package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"

	"./manager"
)

type MarkdownsHandler struct {
	markdownsManeger *manager.MarkdownsManeger
	templatesPath    string
}

func (h *MarkdownsHandler) Index(w http.ResponseWriter, req *http.Request) {
	t := template.New("index.html")
	if h.templatesPath == "" {
		t, _ = t.Parse(DefaultTemplateBox.String("index.html"))
	} else {
		t, _ = t.ParseFiles(fmt.Sprintf("%s%sindex.html", h.templatesPath, string(os.PathSeparator)))
	}
	t.Execute(w, h.markdownsManeger.GetFileList())
}

func (h *MarkdownsHandler) ReaderHander(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName, ok := vars["filename"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := h.markdownsManeger.GetFile(fileName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(body)
	}
}
func (h *MarkdownsHandler) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	if h.markdownsManeger.Reflesh() {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("under refleshing"))
	}

}

var (
	hander                                   *MarkdownsHandler
	TemplateFloder, MarkdownFloder, HttpPort string
	DefaultTemplateBox                       packr.Box
)

func init() {
	flag.StringVar(&HttpPort, "p", "8000", "service port")
	flag.StringVar(&TemplateFloder, "t", "", "path for the floder containning custom html file,not required")
	flag.StringVar(&MarkdownFloder, "f", "markdowns", "the markdows files floder, default: markdows")
	flag.Parse()

	m := manager.New(MarkdownFloder)
	hander = &MarkdownsHandler{
		&m,
		TemplateFloder,
	}
	DefaultTemplateBox = packr.NewBox("templates")
	hander.markdownsManeger.Reflesh()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", hander.Index)
	r.HandleFunc("/file/{filename}/", hander.ReaderHander)
	r.HandleFunc("/update", hander.UpdateHandler)

	http.ListenAndServe(fmt.Sprintf(":%s", HttpPort), r)
}

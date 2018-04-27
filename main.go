package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"zhnxin/markdownServer/manager"
)

type MarkdownsHandler struct {
	markdownsManeger *manager.MarkdownsManeger
	templatesPath    string
}

func (h *MarkdownsHandler) Index(w http.ResponseWriter, req *http.Request) {
	t := template.New("index.html")
	t, _ = t.ParseFiles(fmt.Sprintf("%s%sindex.html", h.templatesPath, string(os.PathSeparator)))
	t.Execute(w, h.markdownsManeger.GetFileList())
}

func (h *MarkdownsHandler) ReaderHander(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
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
	hander *MarkdownsHandler
)

func init() {
	m := manager.New("markdowns")
	hander = &MarkdownsHandler{
		&m,
		"templates",
	}
	hander.markdownsManeger.Reflesh()
}

func main() {
	http.HandleFunc("/", hander.Index)
	http.HandleFunc("/read", hander.ReaderHander)
	http.HandleFunc("/update", hander.UpdateHandler)

	http.ListenAndServe(":8000", nil)
}

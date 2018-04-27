package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

var (
	MARK_DOWNS_PATH = "markdowns"
	TEMPLATES_PATH  = "templates"
	fileMap         map[string]string
	globalLock      *sync.RWMutex
)

func init() {
	globalLock = new(sync.RWMutex)
	updateMarkdownFileList()
}

func updateMarkdownFileList() {
	globalLock.Lock()
	defer globalLock.Unlock()
	fileMap = make(map[string]string)
	files, _ := filepath.Glob(fmt.Sprintf("%s/*", MARK_DOWNS_PATH))
	for _, f := range files {
		fileName := f[strings.LastIndex(f, string(os.PathSeparator))+1 : len(f)-3]
		fileMap[fileName] = f
	}
}

func readMarkdownFile(fileName string) ([]byte, error) {
	globalLock.RLock()
	defer globalLock.RUnlock()
	markdownFile, ok := fileMap[fileName]
	if !ok {
		return nil, fmt.Errorf("file(%s)NotExist", fileName)
	}

	if fileread, err := ioutil.ReadFile(markdownFile); err == nil {
		unsafe := blackfriday.Run(fileread)
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
		return html, nil
	} else {
		return nil, fmt.Errorf("file(%s)ReadFail", fileName)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html")
	t, _ = t.ParseFiles(fmt.Sprintf("%s%sindex.html", TEMPLATES_PATH, string(os.PathSeparator)))
	t.Execute(w, fileMap)
}

func readerHander(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	body, err := readMarkdownFile(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(body)
	}
}
func updateHandler(w http.ResponseWriter, req *http.Request) {
	updateMarkdownFileList()
	w.Write([]byte("success"))
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/read", readerHander)
	http.HandleFunc("/update", updateHandler)

	http.ListenAndServe(":8000", nil)
}

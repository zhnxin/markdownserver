package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/LK4D4/trylock"
	"github.com/microcosm-cc/bluemonday"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// MarkdownFile :markdown file and html byte code
type MarkdownFile struct {
	path string
	data []byte
}

// MarkdownsManeger to update and get markdown file
type MarkdownsManeger struct {
	rwLock       *sync.RWMutex
	updateLock   *trylock.Mutex
	data         map[string]MarkdownFile
	markdownPath string
}

func (m *MarkdownFile) isExist() bool {
	return m.data != nil
}

func (s *MarkdownsManeger) Reflesh() bool {
	if s.updateLock.TryLock() {
		s.rwLock.Lock()
		defer s.updateLock.Unlock()
		defer s.rwLock.Unlock()
		s.data = make(map[string]MarkdownFile)
		files, _ := filepath.Glob(fmt.Sprintf("./%s/*", s.markdownPath))
		for _, f := range files {
			fileName := f[strings.LastIndex(f, string(os.PathSeparator))+1 : len(f)-3]
			s.data[fileName] = MarkdownFile{f, nil}
		}
		return true

	}
	return false
}

func (s *MarkdownsManeger) readMarkdown(fileName string) ([]byte, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	markdownFile, ok := s.data[fileName]
	if !ok {
		return nil, fmt.Errorf("file(%s)NotExist", fileName)
	}
	if markdownFile.isExist() {
		return markdownFile.data, nil
	}

	if fileread, err := ioutil.ReadFile(markdownFile.path); err == nil {
		unsafe := blackfriday.Run(fileread)
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
		markdownFile.data = html
		return html, nil
	}

	return nil, fmt.Errorf("file(%s)ReadFail", fileName)
}

func (s *MarkdownsManeger) GetFileList() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *MarkdownsManeger) GetFile(fileName string) ([]byte, error) {
	return s.readMarkdown(fileName)
}

//New MarkdownsManeger
func New(markdownPath string) MarkdownsManeger {
	return MarkdownsManeger{
		new(sync.RWMutex),
		new(trylock.Mutex),
		map[string]MarkdownFile{},
		markdownPath,
	}
}

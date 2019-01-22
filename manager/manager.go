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
	path       string
	data       []byte
	updateLock *trylock.Mutex
}

// MarkdownsManeger to update and get markdown file
type MarkdownsManeger struct {
	rwLock       *sync.RWMutex
	updateLock   *trylock.Mutex
	data         map[string]*MarkdownFile
	markdownPath string
}

func (m *MarkdownFile) isExist() bool {
	return m.data != nil
}

func (m *MarkdownFile) readMarkdown() ([]byte, error) {
	if m.updateLock.TryLock() {
		defer m.updateLock.Unlock()
		if fileread, err := ioutil.ReadFile(m.path); err == nil {
			unsafe := blackfriday.Run(fileread, blackfriday.WithExtensions(blackfriday.CommonExtensions))
			html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
			m.data = html
			return html, nil
		} else {
			return nil, fmt.Errorf("file(%s)ReadFail", m.path)
		}
	} else {
		m.updateLock.Lock()
		m.updateLock.Unlock()
		if m.isExist() {
			return m.data, nil
		} else {
			return nil, fmt.Errorf("file(%s)ReadFail", m.path)
		}
	}

}

func (s *MarkdownsManeger) Reflesh() bool {
	if s.updateLock.TryLock() {
		s.rwLock.Lock()
		defer s.updateLock.Unlock()
		defer s.rwLock.Unlock()
		s.data = make(map[string]*MarkdownFile)
		files, _ := filepath.Glob(fmt.Sprintf("./%s/*", s.markdownPath))
		for _, f := range files {
			if strings.HasSuffix(f, ".md") || strings.HasSuffix(f, ".MD") {
				fileName := f[strings.LastIndex(f, string(os.PathSeparator))+1 : len(f)-3]
				s.data[fileName] = &MarkdownFile{f, nil, new(trylock.Mutex)}
			}
		}
		return true

	}
	return false
}

func (s *MarkdownsManeger) GetFileList() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *MarkdownsManeger) GetFile(fileName string) ([]byte, error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	markdownFile, ok := s.data[fileName]
	if !ok {
		return nil, fmt.Errorf("file(%s)NotExist", fileName)
	}
	if markdownFile.isExist() {
		return markdownFile.data, nil
	}
	return markdownFile.readMarkdown()
}

//New MarkdownsManeger
func New(markdownPath string) MarkdownsManeger {
	return MarkdownsManeger{
		new(sync.RWMutex),
		new(trylock.Mutex),
		map[string]*MarkdownFile{},
		markdownPath,
	}
}

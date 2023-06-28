package msg

import (
	"fmt"
	"github.com/yonisaka/go-boilerplate/pkg/file"
	"net/http"
	"strings"
	"sync"
)

var msgs map[int]*Message
var once sync.Once

// MessageConfig as messages configuration
type MessageConfig struct {
	Messages []*Message `yaml:"messages"`
}

// Message configuration structure
type Message struct {
	Code     int        `yaml:"code"`
	Contents []*Content `yaml:"contents"`
	contents map[string]*Content
}

// doMap create content map from slice
func (m *Message) doMap() *Message {
	m.contents = make(map[string]*Content, 0)
	for _, c := range m.Contents {
		l := strings.ToLower(c.Lang)
		if _, ok := m.contents[l]; !ok {
			m.contents[l] = c
			continue
		}
	}
	return m
}

// Content message content configuration structure
type Content struct {
	Lang string `yaml:"lang"`
	Text string `yaml:"text"`
}

// Setup initializes messages  from yaml file
// args:
//
//	path: path of message list definition file
//
// returns:
//
//	err: operation error
func Setup(fname string, paths ...string) (err error) {
	var mcfg MessageConfig
	once.Do(func() {
		msgs = make(map[int]*Message, 0)
		for _, p := range paths {
			f := fmt.Sprint(p, fname)
			err = file.ReadFromYAML(f, &mcfg)
			if err != nil {
				continue
			}
			err = nil
		}
	})

	if err != nil {
		err = fmt.Errorf("unable to read config from files %s", err.Error())
		return
	}
	for _, m := range mcfg.Messages {
		if _, ok := msgs[m.Code]; !ok {
			m := &Message{Code: m.Code, Contents: m.Contents}
			msgs[m.Code] = m.doMap()
		}
	}
	return
}

// Get messages by language
func Get(code int, lang string) (text string) {
	lang = cleanLangStr(lang)
	if m, ok := msgs[code]; ok {
		if c, ok := m.contents[lang]; ok {
			text = c.Text
			return
		}
	}
	return
}

// GetCode messages by language
func GetCode(code int) int {
	if m, ok := msgs[code]; ok {
		return m.Code
	}
	return http.StatusUnprocessableEntity
}

// GetMessageCode messages by language
func GetMessageCode(key int, lang string) (code int, text string) {
	lang = cleanLangStr(lang)
	if m, ok := msgs[key]; ok {
		code = m.Code
		if c, ok := m.contents[lang]; ok {
			text = c.Text
			return
		}
	}

	code = http.StatusUnprocessableEntity
	return
}

func cleanLangStr(s string) string {
	return strings.ToLower(strings.Trim(s, " "))
}

// GetAvailableLang func check language
func GetAvailableLang(key int, lang string) bool {
	if m, ok := msgs[key]; ok {
		if _, ok := m.contents[lang]; ok {
			return true
		}
	}

	return false
}

package locales

import (
	"bytes"
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
	"text/template"
)

type Translator interface {
	Translate(ctx context.Context, key string, params map[string]interface{}) string
}

type translator struct {
	defaultLang string
	data        map[string]string
	langs       []string
}

func (t *translator) loadLocalesFile(filepath string) (map[string]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	any := make(map[string]interface{})
	err = yaml.Unmarshal(buf.Bytes(), &any)
	if err != nil {
		return nil, err
	}

	l := langFromFilename(f.Name())
	t.langs = append(t.langs, l)

	return withLanguageName(l, flatten(any)), nil
}

func langFromFilename(name string) string {
	pathParts := strings.Split(name, "/")
	lang := strings.Split(pathParts[len(pathParts)-1], ".")[0]
	return lang
}

func withLanguageName(lang string, flattened map[string]string) map[string]string {
	fl := make(map[string]string)

	for k, v := range flattened {
		fl[fmt.Sprintf("%s.%s", lang, k)] = v
	}

	return fl
}

func concatMap(root map[string]string, maps ...map[string]string) map[string]string {
	for _, mp := range maps {
		for k, v := range mp {
			root[k] = v
		}
	}

	return root
}

func flatten(m map[string]interface{}) map[string]string {
	o := make(map[string]string)
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		default:
			o[k] = fmt.Sprintf("%s", v)
		}
	}
	return o
}

func NewTranslator(defaultLang string, filepath ...string) Translator {
	root := make(map[string]string)

	t := translator{}
	for _, path := range filepath {
		data, err := t.loadLocalesFile(path)
		if err != nil {
			continue
		}

		concatMap(root, data)
	}

	t.data = root
	t.defaultLang = defaultLang
	return &t
}

func (t *translator) Translate(ctx context.Context, key string, params map[string]interface{}) string {
	lang, ok := LanguageFromContext(ctx)
	if !ok {
		lang = t.defaultLang
	}

	if !t.hasLang(lang) {
		lang = t.defaultLang
	}

	text, ok := t.data[fmt.Sprintf("%s.%s", lang, key)]
	if !ok {
		return fmt.Sprintf("invalid key : %s", key)
	}

	tpl, err := template.New(key).Parse(text)
	if err != nil {
		return fmt.Sprintf("invalid key : %s", key)
	}

	str := strings.Builder{}

	err = tpl.Execute(&str, params)
	if err != nil {
		return fmt.Sprintf("invalid key : %s", key)
	}

	return str.String()
}

func (t *translator) hasLang(lang string) bool {
	for _, l := range t.langs {
		if l == lang {
			return true
		}
	}

	return false
}

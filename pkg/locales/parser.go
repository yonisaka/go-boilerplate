package locales

import (
	"golang.org/x/text/language"
	"strings"
)

func ParseIOSLang(in string, def string) string {
	in = strings.ToLower(in)
	switch in {
	case "id", "en":
		return in
	default:
		// size of t == q
		t, _, err := language.ParseAcceptLanguage(in)
		if err != nil {
			return def
		}

		if len(t) == 0 {
			return def
		}

		out := t[0].String()
		if strings.Contains(out, "-") {
			outParts := strings.Split(out, "-")
			out = outParts[0]
		}

		return out
	}
}

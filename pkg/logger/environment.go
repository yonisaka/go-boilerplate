package logger

import "strings"

var envs = map[string]string{
	"production":  "production",
	"staging":     "staging",
	"development": "development",
	"prod":        "production",
	"stg":         "staging",
	"dev":         "development",
	"prd":         "production",
	"green":       "green",
	"blue":        "blue",
}

func Environment(env string) string {
	v, ok := envs[strings.ToLower(strings.Trim(env, " "))]

	if !ok {
		return ""
	}

	return v
}

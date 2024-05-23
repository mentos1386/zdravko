package script

import (
	"html"
	"regexp"
	"strings"
)

func EscapeString(s string) string {
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(html.EscapeString(s), `\n`)
}

func UnescapeString(s string) string {
	return html.UnescapeString(strings.Replace(s, `\n`, "\n", -1))
}

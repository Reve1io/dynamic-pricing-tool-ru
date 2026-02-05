package utils

import (
	"bytes"
	"regexp"
)

var (
	// "key": → уже ок
	// key:   → нужно чинить
	reUnquotedKey = regexp.MustCompile(`([{,]\s*)([a-zA-Z0-9_]+)\s*:`)

	// part: ADUM3471ARSZ  → "ADUM3471ARSZ"
	reUnquotedString = regexp.MustCompile(`:\s*([A-Za-z0-9\-_\.\/]+)([\s,}\]])`)
)

func SanitizeJSON(b []byte) []byte {
	s := string(b)

	// 1. Убираем BOM и мусор
	s = bytes.NewBufferString(s).String()

	// 2. Чиним незакавыченные ключи
	s = reUnquotedKey.ReplaceAllString(s, `$1"$2":`)

	// 3. Чиним незакавыченные строковые значения
	s = reUnquotedString.ReplaceAllString(s, `: "$1"$2`)

	// 4. Убираем висячие запятые
	s = regexp.MustCompile(`,(\s*[}\]])`).ReplaceAllString(s, `$1`)

	return []byte(s)
}

package postgres

import (
	"strconv"
	"strings"
)

// `?, ?, ?, ?` => `$1, $2, $3, $4`
func rebind(query string) string {
	count := 1
	sb := strings.Builder{}
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			sb.WriteString("$")
			sb.WriteString(strconv.Itoa(count))
			count++
			continue
		}
		sb.WriteByte(query[i])
	}

	return sb.String()
}

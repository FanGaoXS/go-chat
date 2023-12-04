package postgres

import (
	"hash/fnv"
	"math"
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

func hashCode(str string) int64 {
	hash := fnv.New64()
	hash.Write([]byte(str))
	hc := hash.Sum64()
	if hc > math.MaxInt64 {
		return int64(hc - math.MaxUint64 - 1)
	}
	return int64(hc)
}

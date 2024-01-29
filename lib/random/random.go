package random

import (
	"math/rand"
	"strings"
	"time"
)

func NewRandomString(length int) string {

	letters := []string{}
	for i := 'a'; i <= 'z'; i++ {
        letters = append(letters, string(i))
    }
	for i := '0'; i <= '9'; i++ {
        letters = append(letters, string(i))
    }

	sb := strings.Builder{}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()));

	for i := 0; i < length; i++ {
		sb.WriteString(letters[rnd.Intn(len(letters))])
	}

	return sb.String()
}
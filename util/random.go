package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generated a random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generates a random string
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// return random owner
func RandomOwner() string {
	return RandomString(6)
}

// random balance
func RandomBalance() int64{
	return RandomInt(0, 1000)
}

// random currency
func RandomCurrency() string{
	currencies := []string{"EUR", "USD", "CAD"}

	n := len(currencies)

	return currencies[rand.Intn(n)]
}

func RandomMoney() int64{
	return RandomInt(0, 10000)
}

func RandomEmail() string{
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

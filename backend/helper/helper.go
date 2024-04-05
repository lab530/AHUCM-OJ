package helper

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm")
	result := make([]byte, n)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func CheckAllEmptyString(str string) bool {
	for _, char := range str {
		if char != ' ' {
			return false
		}
	}
	return true
}

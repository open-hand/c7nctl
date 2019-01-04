package utils

import (
	"math/rand"
	"time"
)

const randomLength = 4

func RandomString(length ...int) string {

	randomLength := randomLength
	if len(length) > 0 {
		randomLength = length[0]
	}
	bytes := make([]byte, randomLength)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < randomLength; i++ {
		bytes[i] = byte(97 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

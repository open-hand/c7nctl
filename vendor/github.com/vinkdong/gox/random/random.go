package random

import (
	"math/rand"
)

const (
	defaultLength = 4
)

type Slice struct {
	Start int
	End   int
}

func Seed(seed int64)  {
	rand.Seed(seed)
}

func RangeInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RangeIntInclude(include ...Slice) int {
	if len(include) < 1 {
		return rand.Int()
	}
	s := rand.Intn(len(include))
	e := include[s]
	return RangeInt(e.Start,e.End)
}

func String(length ...int) string {
	randomLength := defaultLength
	if len(length) > 0 {
		randomLength = length[0]
	}
	
	bytes := make([]byte, randomLength)
	
	for i := 0; i < randomLength; i++ {
		bytes[i] = byte(97 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

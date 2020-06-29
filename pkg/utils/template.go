package utils

import (
	"github.com/chr4/pwgen"
	"github.com/vinkdong/gox/random"
	"math/rand"
	"text/template"
	"time"
)

var C7nFunc = template.FuncMap{
	"randomToken":        randomToken,
	"getImageRepo":       getImageRepo,
	"randomLowCaseToken": randomLowCaseToken,
	"generateAlphaNum":   generateAlphaNum,
}

func getImageRepo(rls string) string {
	return "registry.cn-shanghai.aliyuncs.com/c7n/" + rls
}

func randomToken(length int) string {
	return GenerateRunnerToken(length)
}

func randomLowCaseToken(length int) string {
	return GenerateRunnerToken(length)
}

func generateAlphaNum(length int) string {
	return pwgen.AlphaNum(length)
}

func RandomToken(length int) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		random.Seed(time.Now().UnixNano())
		op := random.RangeIntInclude(random.Slice{Start: 48, End: 57},
			random.Slice{Start: 65, End: 90}, random.Slice{Start: 97, End: 122})
		b[i] = byte(op) //A=65 and Z = 65+25
	}
	return string(b)
}

func GenerateRunnerToken(length int) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		random.Seed(time.Now().UnixNano())
		op := random.RangeIntInclude(random.Slice{Start: 48, End: 57},
			random.Slice{Start: 97, End: 122})
		b[i] = byte(op) //A=65 and Z = 65+25
	}
	return string(b)
}

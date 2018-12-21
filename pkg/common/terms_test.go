package common

import "testing"

//func TestAskAgreeTerms(t *testing.T){
//	AskAgreeTerms()
//}

func TestCheckMatch(t *testing.T)  {
	input := Input{
		Regex: "^\\S{8}$",
		Include: []KV{KV{Name:"必须包含数字",Value:"\\d+"},KV{Name:"必须包含大写",Value:"[A-Z]+"}},
	}

	if CheckMatch("abcdef",input) {
		t.Error("regex check error")
	}

	if !CheckMatch("ce123Qwe",input) {
		t.Error("regex check error")
	}

	if CheckMatch("ce123456",input) {
		t.Error("regex check error")
	}

	if CheckMatch("cdsfaQes",input) {
		t.Error("regex check error")
	}
}
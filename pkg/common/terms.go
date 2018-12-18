package common

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"regexp"
	"github.com/vinkdong/gox/log"
)

type Input struct {
	Enabled  bool
	Regex    string
	Tip      string
	Password bool
	Include []KV
	Exclude []KV
}

type KV struct {
	Name  string
	Value string
}

func AskAgreeTerms() {
	input := Input{
		Password: false,
		Tip:      "为了提高用户体验，程序会收集一些非敏感信息上传到我们服务器，具体包括:主机内存大小、CPU数量/频率、IP(仅用于确认您所在的地区不会进行储存)、Kubernetes版本等\nIn order to improve the user experience, the program will collect some non-sensitive information to upload to our server, including: host memory size, CPU frequency, ip(just for ensure where your city), Kubernetes version, etc. \n同意请输入Y，不同意请输入N。\nagree to enter Y, do not agree to enter N. [Y/N]:   ",
		Regex:    "^(y|Y|n|N)$",
	}
	r, err := AcceptUserInput(input)
	if err != nil {
		log.Error(err)
		os.Exit(127)
	}
	if r == "n" || r == "N" {
		os.Exit(151)
	}
}

func AcceptUserInput(input Input) (string, error) {
start:
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(input.Tip)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "",err
	}
	text = strings.Trim(text,"\n")

	if !CheckMatch(text,input) {
		goto start
	}
	return text, nil
}

func CheckMatch(value string, input Input) bool {

	r := regexp.MustCompile(input.Regex)
	if !r.MatchString(value) {
		log.Errorf("输入不满足需求")
		return false
	}

	for _, include := range input.Include {
		r := regexp.MustCompile(include.Value)
		if !r.MatchString(value) {
			log.Errorf(include.Name)
			return false
		}
	}

	for _,exclude := range input.Exclude {
		r := regexp.MustCompile(exclude.Value)
		if r.MatchString(value) {
			log.Errorf(exclude.Name)
			return false
		}
	}

	return true
}
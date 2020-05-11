package utils

import (
	"bufio"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/hashicorp/go-version"
	"github.com/vinkdong/gox/log"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"strings"
	"syscall"
)

func AskAgreeTerms() {
	input := context.Input{
		Password: false,
		Tip:      "为了提高用户体验，程序会收集一些非敏感信息上传到我们服务器，具体包括:主机内存大小、CPU数量/频率、Kubernetes版本\nIn order to improve the user experience, the program will collect some non-sensitive information to upload to our server, including: host memory size, CPU frequency, Kubernetes version. \n同意请输入Y，不同意请输入N。\nagree to enter Y, do not agree to enter N. [Y/N]:   ",
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

func AcceptUserInput(input context.Input) (string, error) {
	if input.Password {
		return AcceptUserPassword(input)
	}
start:
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(input.Tip)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.Trim(text, "\n")

	if !CheckMatch(text, input) {
		goto start
	}
	return text, nil
}

func AcceptUserPassword(input context.Input) (string, error) {
start:
	fmt.Print(input.Tip)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	if !CheckMatch(string(bytePassword[:]), input) {
		goto start
	}

	if !input.Twice {
		return string(bytePassword[:]), nil
	}

	fmt.Print("请再输入一次:")
	bytePassword2, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	if len(bytePassword2) != len(bytePassword) {
		log.Error("两次输入长度不符")
		goto start
	}
	for k, v := range bytePassword {
		if bytePassword2[k] != v {
			log.Error("两次输入不同")
			goto start
		}
	}

	fmt.Println("waiting...")
	return string(bytePassword[:]), nil
}

func CheckMatch(value string, input context.Input) bool {

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

	for _, exclude := range input.Exclude {
		r := regexp.MustCompile(exclude.Value)
		if r.MatchString(value) {
			log.Errorf(exclude.Name)
			return false
		}
	}

	return true
}

func CheckVersion(versionRaw, constraint string) (bool, error) {
	v1, err := version.NewVersion(versionRaw)
	if err != nil {
		return false, err
	}
	constraints, err := version.NewConstraint(constraint)
	if err != nil {
		return false, err
	}
	return constraints.Check(v1), nil
}

func ConditionSkip() bool {
	//todo skip some test in conditions
	return true
}

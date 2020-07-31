package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/vinkdong/gox/log"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"regexp"
	"strings"
)

type Input struct {
	Enabled  bool
	Regex    string
	Tip      string
	Password bool
	Include  []KV
	Exclude  []KV
	Twice    bool
}

type KV struct {
	Name  string
	Value string
}

func AskAgreeTerms() {
	input := Input{
		Password: false,
		Tip:      "为了提高用户体验，程序会收集一些非敏感信息上传到我们服务器，具体包括:主机内存大小、CPU数量/频率、Kubernetes版本\nIn order to improve the user experience, the program will collect some non-sensitive information to upload to our server, including: host memory size, CPU frequency, Kubernetes version. \n同意请输入Y，不同意请输入N。\nagree to enter Y, do not agree to enter N. [Y/N]: ",
		Regex:    "^(y|Y|n|N)$",
	}
	r, err := AcceptUserInput(input)
	if err != nil {
		log.Error(err)
		os.Exit(127)
	}
	// 拒绝协议，直接退出
	if r == "n" || r == "N" {
		os.Exit(151)
	}
}

func AcceptUserInput(input Input) (string, error) {
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
	text = strings.Trim(text, "\r")

	if !CheckMatch(text, input) {
		goto start
	}
	return text, nil
}

func AcceptUserPassword(input Input) (string, error) {
start:
	fmt.Print(input.Tip)
	// TODO
	bytePassword, err := readPassword("") // terminal.ReadPassword(int(os.Stdin.Fd()))
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
	bytePassword2, err := readPassword("") // terminal.ReadPassword(int(os.Stdin.Fd()))
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

	return string(bytePassword[:]), nil
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

func readPassword(prompt string) (pw []byte, err error) {
	fd := int(os.Stdin.Fd())
	if terminal.IsTerminal(fd) {
		fmt.Fprint(os.Stderr, prompt)
		pw, err = terminal.ReadPassword(fd)
		fmt.Fprintln(os.Stderr)
		return
	}

	var b [1]byte
	for {
		n, err := os.Stdin.Read(b[:])
		// terminal.ReadPassword discards any '\r', so we do the same
		if n > 0 && b[0] != '\r' {
			if b[0] == '\n' {
				return pw, nil
			}
			pw = append(pw, b[0])
			// limit size, so that a wrong input won't fill up the memory
			if len(pw) > 1024 {
				err = errors.New("password too long")
			}
		}
		if err != nil {
			// terminal.ReadPassword accepts EOF-terminated passwords
			// if non-empty, so we do the same
			if err == io.EOF && len(pw) > 0 {
				err = nil
			}
			return pw, err
		}
	}
}

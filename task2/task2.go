package task2

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	multipleDigitRe = regexp.MustCompile(`(^|[^\\])(\\\\)*\d\d`)
	commandRe       = regexp.MustCompile(`(^|[^\\])(\\\\)*\d`)
	slashRe         = regexp.MustCompile(`\\\\?`)
)

func StringUnpack(src string) string {
	if multipleDigitRe.MatchString(src) {
		return ""
	}
	res := commandRe.ReplaceAllStringFunc(src, unpackCommand)
	res = slashRe.ReplaceAllStringFunc(res, handleSlash)
	return res
}

func unpackCommand(cmd string) string {
	cmdR := []rune(cmd)
	cmdRLen := len(cmdR)
	if cmdRLen == 1 {
		return ""
	}
	cnt, _ := strconv.Atoi(string(cmdR[cmdRLen-1:])) // ignore error, because of regexp guarantee
	repeatR := cmdR[:cmdRLen-1]
	prefixR := []rune(``)
	if cmdRLen > 2 { // if slashes
		prefixR = repeatR[:len(repeatR)-2]
		repeatR = []rune(`\\`)
	}
	return string(prefixR) + strings.Repeat(string(repeatR), cnt)
}

func handleSlash(src string) string {
	if src == `\\` {
		return `\`
	}
	return ``
}

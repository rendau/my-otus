package task2

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimple(t *testing.T) {
	var src = "a4bc2d5e"
	var dst = "aaaabccddddde"

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestWithoutCommand(t *testing.T) {
	var src = "abcd"
	var dst = "abcd"

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestCyrillic(t *testing.T) {
	var src = "ы4вж2ё5ь"
	var dst = "ыыыывжжёёёёёь"

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestStartsWithDigit(t *testing.T) {
	var src = "4qwe"
	var dst = "qwe"

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestMultipleDigits(t *testing.T) {
	var src = "45"
	var dst = ""

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestEscapedSymbolsNotInCommand(t *testing.T) {
	var src = `\q\\\we4`
	var dst = `q\weeee`

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestEscapedDigits(t *testing.T) {
	var src = `qwe\4\5`
	var dst = "qwe45"

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestEscapedDigitSymbol(t *testing.T) {
	var src = `qwe\45`
	var dst = "qwe44444"

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestMultipleDigitsWithEscapedSlashPrefix(t *testing.T) {
	var src = `qwe\\45`
	var dst = ""

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestEscapedSlashSymbol(t *testing.T) {
	var src = `qwe\\5`
	var dst = `qwe\\\\\`

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestTripleSlash1(t *testing.T) {
	var src = `qwe\\\4\\\5`
	var dst = `qwe\4\5`

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestTripleSlash2(t *testing.T) {
	var src = `qwe\\\5`
	var dst = `qwe\5`

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestTripleSlash3(t *testing.T) {
	var src = `qwe\\\\5`
	var dst = `qwe\\\\\\`

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

func TestTripleSlash4(t *testing.T) {
	var src = `qwe\\\45`
	var dst = `qwe\44444`

	require.Equal(t, dst, StringUnpack(src), "case '"+src+"'")
}

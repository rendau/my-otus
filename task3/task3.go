package task2

import (
	"regexp"
	"sort"
)

var (
	wsRe = regexp.MustCompile(`\s+`)
)

type wordInfoSt struct {
	v   string
	cnt uint64
}

/*
	Здесь:
		src - исходный текст
		n - нужное количество слов для возврата (по заданию это цифра 10)
*/
func Fun1(src string, n int) (result []string) {
	words := wsRe.Split(src, -1)
	sort.Strings(words)
	wordInfoList := make([]*wordInfoSt, 0, len(words))
	var wordInfo *wordInfoSt
	for _, w := range words {
		if w == "" {
			continue
		}
		if wordInfo == nil || w != wordInfo.v {
			wordInfo = &wordInfoSt{
				v:   w,
				cnt: 1,
			}
			wordInfoList = append(wordInfoList, wordInfo)
		} else {
			wordInfo.cnt += 1
		}
	}
	sort.Slice(wordInfoList, func(i, j int) bool {
		if wordInfoList[i].cnt == wordInfoList[j].cnt {
			return wordInfoList[i].v < wordInfoList[j].v
		}
		return wordInfoList[i].cnt > wordInfoList[j].cnt
	})
	for i := 0; i < len(wordInfoList) && i < n; i++ {
		result = append(result, wordInfoList[i].v)
	}
	return result
}

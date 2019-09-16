package task2

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type caseSt struct {
	s string
	r []string
}

func TestCommon(t *testing.T) {
	cases := []caseSt{
		{
			s: "www  dsa  www     qwe rty \nwww rrr dsa",
			r: []string{"dsa", "www"},
		},
		{
			s: "      \n      www  dsa  www  \n\n\n\n   qwe rty \nwww rrr dsa     \n\n      ",
			r: []string{"dsa", "www"},
		},
	}

	for _, c := range cases {
		result := Fun1(c.s, 2)
		require.True(t, isSameStringSlices(result, c.r), fmt.Sprintf(`case %q, right: %v, actual: %v`, c.s, c.r, result))
	}
}

/*
	эта функция сравнивает два слайса по значениям, даже если значения не в одинаковом порядке
*/
func isSameStringSlices(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for _, x := range s1 {
		found := false
		for _, y := range s2 {
			if x == y {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

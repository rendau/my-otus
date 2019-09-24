package task4

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPushFrontToEmptyList(t *testing.T) {
	var l = List{}
	var v1 = 5

	i1 := l.PushFront(v1)

	require.NotNil(t, i1)
	require.Equal(t, v1, i1.Value())
	require.Equal(t, uint64(1), l.Len())
	require.Equal(t, l.First(), i1)
	require.Equal(t, l.First(), l.Last())
	require.Nil(t, i1.Prev())
	require.Nil(t, i1.Next())
}

func TestPushFrontToNotEmptyList(t *testing.T) {
	var l = List{}
	var v1 = 5
	var v2 = 7

	i1 := l.PushFront(v1)
	i2 := l.PushFront(v2)

	require.NotNil(t, i1)
	require.NotNil(t, i2)
	require.NotEqual(t, i1, i2)
	require.Equal(t, v1, i1.Value())
	require.Equal(t, v2, i2.Value())
	require.Equal(t, uint64(2), l.Len())
	require.Equal(t, l.First(), i2)
	require.Equal(t, l.Last(), i1)
	require.Nil(t, i2.Prev())
	require.Equal(t, i2.Next(), i1)
	require.Nil(t, i1.Next())
	require.Equal(t, i1.Prev(), i2)
}

func TestPushBackToEmptyList(t *testing.T) {
	var l = List{}
	var v1 = 5

	i1 := l.PushBack(v1)

	require.NotNil(t, i1)
	require.Equal(t, v1, i1.Value())
	require.Equal(t, uint64(1), l.Len())
	require.Equal(t, l.First(), i1)
	require.Equal(t, l.First(), l.Last())
	require.Nil(t, i1.Prev())
	require.Nil(t, i1.Next())
}

func TestPushBackToNotEmptyList(t *testing.T) {
	var l = List{}
	var v1 = 5
	var v2 = 7

	i1 := l.PushBack(v1)
	i2 := l.PushBack(v2)

	require.NotNil(t, i1)
	require.NotNil(t, i2)
	require.NotEqual(t, i1, i2)
	require.Equal(t, v1, i1.Value())
	require.Equal(t, v2, i2.Value())
	require.Equal(t, uint64(2), l.Len())
	require.Equal(t, l.First(), i1)
	require.Equal(t, l.Last(), i2)
	require.Nil(t, i1.Prev())
	require.Equal(t, i1.Next(), i2)
	require.Nil(t, i2.Next())
	require.Equal(t, i2.Prev(), i1)
}

func TestRemove(t *testing.T) {
	var l = List{}
	var v1 = 5
	var v2 = 7
	var v3 = 9

	i1 := l.PushBack(v1)
	i2 := l.PushBack(v2)
	i3 := l.PushBack(v3)

	l.Remove(i3)

	require.Equal(t, uint64(2), l.Len())
	require.Equal(t, l.First(), i1)
	require.Equal(t, l.Last(), i2)
	require.Nil(t, i3.Prev())
	require.Nil(t, i3.Next())

	l.Remove(i3)

	require.Equal(t, uint64(2), l.Len())
	require.Equal(t, l.First(), i1)
	require.Equal(t, l.Last(), i2)
	require.Nil(t, i3.Prev())
	require.Nil(t, i3.Next())

	l.Remove(i1)

	require.Equal(t, uint64(1), l.Len())
	require.Equal(t, l.First(), i2)
	require.Equal(t, l.Last(), i2)
	require.Nil(t, i1.Prev())
	require.Nil(t, i1.Next())

	l.Remove(i2)

	require.Equal(t, uint64(0), l.Len())
	require.Nil(t, l.First())
	require.Nil(t, l.Last())
}

package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	testSrcFPath = "test_src.txt"
	testSrcData  = "123456789"
	testDstFPath = "test_dst.txt"
)

func TestMain(m *testing.M) {
	err := ioutil.WriteFile(testSrcFPath, []byte(testSrcData), 0666)
	if err != nil {
		log.Fatalln("Fail to create " + testSrcFPath + " file")
	}

	code := m.Run()

	_ = os.Remove(testSrcFPath)
	_ = os.Remove(testDstFPath)

	os.Exit(code)
}

func TestGeneral(t *testing.T) {
	err := fileCopy(testSrcFPath, testDstFPath, 0, 0)
	require.Nil(t, err, err)

	require.Equal(t, testSrcData, readDstFile())
}

func TestOffset(t *testing.T) {
	err := fileCopy(testSrcFPath, testDstFPath, 3, 0)
	require.Nil(t, err, err)

	require.Equal(t, testSrcData[3:], readDstFile())

	err = fileCopy(testSrcFPath, testDstFPath, 1000, 0)
	require.Nil(t, err, err)

	require.Equal(t, "", readDstFile())

	err = fileCopy(testSrcFPath, testDstFPath, -1000, 0)
	require.Nil(t, err, err)

	require.Equal(t, testSrcData, readDstFile())
}

func TestLimit(t *testing.T) {
	err := fileCopy(testSrcFPath, testDstFPath, 0, 2)
	require.Nil(t, err, err)

	require.Equal(t, testSrcData[:2], readDstFile())

	err = fileCopy(testSrcFPath, testDstFPath, 0, 1000)
	require.Nil(t, err, err)

	require.Equal(t, testSrcData, readDstFile())

	err = fileCopy(testSrcFPath, testDstFPath, 0, -1000)
	require.Nil(t, err, err)

	require.Equal(t, testSrcData, readDstFile())

	err = fileCopy(testSrcFPath, testDstFPath, 5, 2)
	require.Nil(t, err, err)

	require.Equal(t, testSrcData[5:7], readDstFile())

	err = fileCopy(testSrcFPath, testDstFPath, 5, 1000)
	require.Nil(t, err, err)

	require.Equal(t, testSrcData[5:], readDstFile())
}

func readDstFile() string {
	d, err := ioutil.ReadFile(testDstFPath)
	if err != nil {
		return ""
	}
	return string(d)
}

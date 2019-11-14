package main

import (
	"bytes"
	"github.com/rendau/my-otus/task7/internal"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testEnvDirPath = "testEnvDir"
)

func TestReadDir(t *testing.T) {
	_ = os.RemoveAll(testEnvDirPath)
	err := os.Mkdir(testEnvDirPath, os.ModePerm)
	if err != nil {
		log.Fatalln("Fail to create " + testEnvDirPath + " dir")
	}
	defer func() { _ = os.RemoveAll(testEnvDirPath) }()

	err = ioutil.WriteFile(filepath.Join(testEnvDirPath, "VAR1"), []byte("val1"), os.ModePerm)
	if err != nil {
		log.Fatalln("Fail to create VAR1 file")
	}

	err = ioutil.WriteFile(filepath.Join(testEnvDirPath, "VAR2"), []byte("val2"), os.ModePerm)
	if err != nil {
		log.Fatalln("Fail to create VAR2 file")
	}

	err = ioutil.WriteFile(filepath.Join(testEnvDirPath, "VAR3"), []byte(""), os.ModePerm)
	if err != nil {
		log.Fatalln("Fail to create VAR3 file")
	}

	err = os.Mkdir(filepath.Join(testEnvDirPath, "sub_dir"), os.ModePerm)
	if err != nil {
		log.Fatalln("Fail to create sub_dir dir")
	}

	err = ioutil.WriteFile(filepath.Join(testEnvDirPath, "sub_dir", "VAR99"), []byte("val99"), os.ModePerm)
	if err != nil {
		log.Fatalln("Fail to create VAR99 file")
	}

	envs, err := internal.ReadDir(testEnvDirPath)
	require.Nil(t, err, err)

	require.Equal(t, 3, len(envs))
	require.Equal(t, "val1", envs["VAR1"])
	require.Equal(t, "val2", envs["VAR2"])
	var3v, ok := envs["VAR3"]
	require.True(t, ok)
	require.Equal(t, "", var3v)
}

func TestRunCmd(t *testing.T) {
	stdOutBuf := new(bytes.Buffer)

	exitCode := internal.RunCmd(
		[]string{"printenv"},
		map[string]string{
			"VAR1": "val1",
			"VAR2": "val2",
			"VAR3": "",
		},
		stdOutBuf,
		nil,
	)
	require.Equal(t, 0, exitCode)

	envs := map[string]string{}
	for _, kv := range strings.Split(stdOutBuf.String(), "\n") {
		kvArr := strings.Split(kv, "=")
		if len(kvArr) == 2 && kvArr[0] != "" {
			envs[kvArr[0]] = kvArr[1]
		}
	}

	require.Equal(t, 3, len(envs))
	require.Equal(t, "val1", envs["VAR1"])
	require.Equal(t, "val2", envs["VAR2"])
	var3v, ok := envs["VAR3"]
	require.True(t, ok)
	require.Equal(t, "", var3v)
}

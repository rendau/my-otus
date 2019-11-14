package internal

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ReadDir(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := map[string]string{}

	var fileContBytes []byte
	var fileContLine string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileContBytes, err = ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		for _, fileContLine = range strings.Split(string(fileContBytes), "\n") {
			res[file.Name()] = strings.TrimSpace(fileContLine)
		}
	}

	return res, nil
}

func RunCmd(cmd []string, envs map[string]string, stdOut, stdErr io.Writer) int {
	if len(cmd) == 0 {
		log.Fatalln("Bad cmd")
	}

	var env []string
	for k, v := range envs {
		env = append(env, k+"="+v)
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = env
	if stdOut != nil {
		command.Stdout = stdOut
	} else {
		command.Stdout = os.Stdout
	}
	if stdErr != nil {
		command.Stderr = stdErr
	} else {
		command.Stderr = os.Stderr
	}

	err := command.Run()
	if err != nil {
		log.Println("Fail to exec command:", err)
		if command.ProcessState != nil {
			return command.ProcessState.ExitCode()
		}
		return 111
	}

	return 0
}

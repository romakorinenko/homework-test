package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			return nil, errors.New("incorrect file name")
		}

		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		var resValue string
		if fileInfo.Size() != 0 {
			value, err := getValueFromFile(path.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			resValue = value
		}

		env[file.Name()] = EnvValue{
			Value:      resValue,
			NeedRemove: fileInfo.Size() == 0,
		}
	}

	return env, nil
}

func getValueFromFile(fullName string) (string, error) {
	fileContent, err := os.Open(fullName)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fileContent.Close()
	}()

	buf := bufio.NewReader(fileContent)
	line, err := buf.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	value := strings.ReplaceAll(string(line), "\x00", "\n")

	return strings.TrimRight(value, " \t\n"), nil
}

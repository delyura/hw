package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
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
		log.Println("Error read dir")
		return nil, err
	}

	env := make(Environment)

	for _, file := range files {
		var value string
		value, err = getValue(dir, file.Name())
		if err != nil {
			return nil, errors.New("invalid argument")
		}

		var infoFile fs.FileInfo
		infoFile, err = file.Info()
		if err != nil {
			return nil, errors.New("error file info")
		}
		env[file.Name()] = EnvValue{
			Value:      value,
			NeedRemove: infoFile.Size() == 0,
		}
	}
	return env, nil
}

func getValue(dir string, fileName string) (string, error) {
	fullPath := path.Join(dir, fileName)
	if strings.Contains(fileName, "=") {
		return "", errors.New("Name of file contain char `=`") //nolint:all
	}
	f, err := os.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed open file")
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("failed close file: %v", err)
		}
	}()

	buffer := bufio.NewReader(f)
	line, err := buffer.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", errors.New("failed to read file")
	}
	val := strings.ReplaceAll(string(line), "\x00", "\n")
	val = strings.TrimRight(val, " \t\n")

	return val, nil
}

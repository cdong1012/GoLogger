package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var fileMutex sync.Mutex

func createFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); err == nil {
		// file already exists
		return os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
	} else if os.IsNotExist(err) {
		// file not exists
		return os.Create(path)
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return nil, errors.New("File invalid")
	}
}

func appendToFile(file *os.File, content []byte) (int, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return file.Write(content)
}

// read file according to length of buffer
// return 0 if EOF
func readFile(file *os.File, buffer []byte) (int, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return file.Read(buffer)
}

func getFileSize(file *os.File) (int64, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	fileInfo, err := file.Stat()
	if err != nil {
		return -1, err
	}
	return fileInfo.Size(), nil
}

func getFileName(file *os.File) (string, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	fileInfo, err := file.Stat()
	if err != nil {
		return "Error", err
	}
	return fileInfo.Name(), nil
}

func copyFile(file *os.File, newPath string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	var fileSize int64
	var err error
	if fileSize, err = getFileSize(file); err != nil {
		return err
	}

	fileBuffer := make([]byte, fileSize)

	_, err = readFile(file, fileBuffer)

	if err != nil {
		return err
	}

	newFile, err := createFile(newPath)
	defer newFile.Close()
	if err != nil {
		return err
	}

	_, err = appendToFile(newFile, fileBuffer)

	if err != nil {
		return err
	}
	return nil
}

func selfReplicate(newPath string) error {
	input, err := ioutil.ReadFile(os.Args[0])
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(newPath, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

func checkError(err error, msg string) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[x] Fatal error: %s\n[xxx] "+msg, err.Error())
		return false
	}
	return true
}

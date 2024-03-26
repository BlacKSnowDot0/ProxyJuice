package Utility

import "os"

func AppendToFile(path string, content string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return
	}
}

func WriteToFile(path string, content string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return
	}
}

package helper

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func ReadFileTxt(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		PrettyLog("error", fmt.Sprintf("Failed to read file txt: %v", err))
		return nil
	}
	defer file.Close()

	var value []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		value = append(value, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		PrettyLog("error", fmt.Sprintf("Error reading file: %v", err))
	}

	return value
}

func CheckFileOrFolder(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func RandomNumber(min int, max int) int {
	return rand.Intn(max-min) + min
}

func FindKeyByValue(m map[string]interface{}, value interface{}) []string {
	var key []string
	for k, v := range m {
		if v == value {
			key = append(key, k)
		}
	}
	return key
}

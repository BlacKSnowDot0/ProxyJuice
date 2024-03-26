package Utility

import (
	"bufio"
	"net/http"
	"strings"
)

func GatherLines(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	var lines []string

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	return lines
}

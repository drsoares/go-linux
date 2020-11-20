package proc

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func extractCmdLine(basePath string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(basePath, cmdLine))
	if err != nil {
		return "", err
	}
	c := 0 // cursor of useful bytes
	for i := 0; i < len(content)-1; i++ {
		// Check if next byte is not a '\0' byte.
		if content[i+1] != byte(0) {
			// Offset must match a '\0' byte.
			c = i + 2
			if content[i] == 0 {
				content[i] = byte(0x20)
			}
		}
	}
	return strings.TrimSpace(string(content[0:c])), nil
}

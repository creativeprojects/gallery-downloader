package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// uniqueName checks the file already exists: if yes it adds a (n) at the end
func uniqueName(filename string) string {
	if _, err := os.Stat(filename); err == nil || os.IsExist(err) {
		extension := path.Ext(filename)
		base := strings.TrimSuffix(filename, extension)
		index := 1
		for {
			filename = fmt.Sprintf("%s(%d)%s", base, index, extension)
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				return filename
			}
			index++
		}
	}
	return filename
}

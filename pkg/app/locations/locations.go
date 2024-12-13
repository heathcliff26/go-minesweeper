package locations

import "os"

const appName = "go-minesweeper"

var saveFolder = ""

func SaveFolder() (string, error) {
	if saveFolder != "" {
		return saveFolder, os.MkdirAll(saveFolder, 0755)
	}
	return os.Getwd()
}

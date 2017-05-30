package mrspider

import (
	"os/exec"
	"os"
	"log"
	"path/filepath"
)

func Max(a, b int) int{
	if a > b {
		return a
	} else {
		return b
	}
}

func GetCurrentDir() string {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	dir := filepath.Dir(path)
	return dir
}

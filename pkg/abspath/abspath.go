package abspath

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// GetAbsolutePath returns absolute path to project root.
func GetAbsolutePath() string {
	var absolutePath string

	if strings.Contains(os.Args[0], "/tmp/") || strings.Contains(os.Args[0], "/var/") {
		absolutePath = getAbsolutePathFromLocalGoRun()
	} else {
		absolutePath = getAbsolutePathFromBinaryRun()
	}

	return absolutePath
}

func getAbsolutePathFromBinaryRun() string {
	p, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}

	absolutePath := filepath.Dir(filepath.Dir(p))

	return absolutePath + "/"
}

func getAbsolutePathFromLocalGoRun() string {
	stackFrame := 2
	_, b, _, _ := runtime.Caller(stackFrame)
	d := path.Join(path.Dir(b) + "/../../")

	return d + "/"
}

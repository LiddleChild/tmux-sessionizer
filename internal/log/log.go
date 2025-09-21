package log

import (
	"flag"
	"os"
)

var (
	DebugFlag = flag.Bool("debug", false, "debug")

	LogFile *os.File
)

func init() {
	flag.Parse()
}

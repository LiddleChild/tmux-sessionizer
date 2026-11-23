package flag

import "flag"

var (
	VersionFlag bool
	DebugFlag   bool
)

func Init() {
	flag.BoolVar(&VersionFlag, "version", false, "version")
	flag.BoolVar(&DebugFlag, "debug", false, "debug")

	flag.Parse()
}

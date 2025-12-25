package utils

import (
	"os"
	"strings"
)

type Args struct {
	Input       string
	Output      []string
	Lang        string
	PackageBase string
	InjectBson  bool
}

func GetArgs() Args {
	args := Args{
		Input:       "minibin-shemas",
		Output:      []string{"minibin-output"},
		Lang:        "go",
		PackageBase: "",
		InjectBson:  false,
	}
	rawArgs := os.Args[1:]
	i := 0
	for i < len(rawArgs) {
		if i+1 >= len(rawArgs) {
			break
		}
		value := rawArgs[i]
		switch value {
		case "-o":
			args.Output = strings.Split(rawArgs[i+1], "/")
		case "-i":
			args.Input = rawArgs[i+1]
		case "-l":
			args.Lang = rawArgs[i+1]
		case "-pkg":
			args.PackageBase = rawArgs[i+1]
		case "-bson":
			args.InjectBson = true
		}
		i += 2
	}
	return args
}

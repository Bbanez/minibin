package utils

import (
	"os"
	"slices"
	"strings"
)

type Args struct {
	StringTemplates bool
	Setup           bool
	Input           []string
	Output          []string
	Lang            string
	PackageBase     string
	InjectBson      bool
	Clear           bool
	ProjectBuild    bool
	CreateRelease   bool
}

type ArgInfo struct {
	flags         []string
	typ           string
	desc          string
	allowedValues *[]string
	parse         func(argv []string, at_idx int, self *ArgInfo, args *Args)
}

var argInfos = []ArgInfo{
	{
		flags: []string{"--input", "-input", "input", "--i", "-i", "i"},
		desc:  "Where to find schema files, Default: minibin-schemas",
		typ:   "string",
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.Input = strings.Split(argv[at_idx], "/")
		},
	},
	{
		flags: []string{"--output", "-output", "output", "--o", "-o", "o"},
		desc:  "Where to output generated files",
		typ:   "string",
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.Output = strings.Split(argv[at_idx], "/")
		},
	},
	{
		flags:         []string{"--language", "-language", "language", "--l", "-l", "l"},
		desc:          "For which language to generate code, Default: go",
		typ:           "string",
		allowedValues: &[]string{"go", "cpp", "rust", "ts"},
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.Lang = argv[at_idx]
		},
	},
	{
		flags: []string{"--package", "-package", "package", "--pkg", "-pkg", "pkg"},
		desc:  "Base package name for generated code, Default: empty",
		typ:   "string",
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.PackageBase = argv[at_idx]
		},
	},
	{
		flags: []string{"--bson", "-bson", "bson"},
		desc:  "Base package name for generated code, Default: empty",
		typ:   "bool",
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.InjectBson = true
		},
	},
	{
		flags: []string{"--clear", "-clear", "clear", "--c", "-c", "c"},
		desc:  "Clear minibin output directory",
		typ:   "bool",
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.Clear = true
		},
	},
	{
		flags: []string{"--build-project", "-build-project", "build-project", "--bp", "-bp", "bp"},
		desc:  "Build the project and create executable binaries",
		typ:   "bool",
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.ProjectBuild = true
		},
	},
	{
		flags: []string{"--release", "-release", "release", "--r", "-r", "r"},
		desc:  "Create a release and deploy",
		typ:   "bool",
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.CreateRelease = true
		},
	},
	{
		flags: []string{"--string-templates", "string-templates", "st"},
		desc:  "Generate string templates",
		typ:   "bool",
		parse: func(argv []string, at_idx int, self *ArgInfo, args *Args) {
			args.StringTemplates = true
		},
	},
}

func GetArgs() *Args {
	argv := os.Args[1:]
	args := &Args{}
	for _, info := range argInfos {
		for idx, arg_value := range argv {
			if !slices.Contains(info.flags, arg_value) {
				continue
			}
			info.parse(argv, idx, &info, args)
		}
	}
	return args
}

package main

import (
	"github.com/bbanez/minibin/src/utils"
	"github.com/bbanez/minibin/tools"
)

func main() {
	args := utils.GetArgs()
	if args.ProjectBuild {
		tools.BuildProject(args)
	} else if args.Clear {
		tools.ClearOutput(args)
	} else if len(args.Input) > 0 && len(args.Output) > 0 && args.Lang != "" {
		tools.Generate(args)
	} else if args.StringTemplates {
		tools.GenerateGoStringTemplates()
	} else {
		args.PrintHelp()
		panic(
			"Invalid arguments combination provided.",
		)
	}
}

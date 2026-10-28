package main

import "github.com/bbanez/minibin/src/utils"

func main() {
	args := utils.GetArgs()
	if args.StringTemplates {
		GenerateGoStringTemplates()
	}
}

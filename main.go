package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bbanez/minibin/src/parser"
	parser_go "github.com/bbanez/minibin/src/parser/go"
	parser_ts "github.com/bbanez/minibin/src/parser/ts"
	"github.com/bbanez/minibin/src/schema"
	"github.com/bbanez/minibin/src/utils"
)

const version = "v0.1.3"

type BuildOutputs struct {
	Name     string
	Platform string
	Arch     string
}

func main() {
	args := utils.GetArgs()
	fs := utils.NewFS(&args.Output)
	if args.ProjectBuild {
		buildOutputs := []BuildOutputs{
			{
				Name:     "Linux",
				Platform: "linux",
				Arch:     "amd64",
			},
			{
				Name:     "Linux",
				Platform: "linux",
				Arch:     "arm64",
			},
			{
				Name:     "Windows",
				Platform: "windows",
				Arch:     "amd64",
			},
			{
				Name:     "MacOS",
				Platform: "darwin",
				Arch:     "arm64",
			},
		}
		for i := range buildOutputs {
			output := buildOutputs[i]
			fmt.Printf(
				"Build %s %s binary...",
				output.Name, output.Arch,
			)
			ext := ""
			if output.Platform == "windows" {
				ext = ".exe"
			}
			cmd := exec.Command(
				"go",
				"build",
				"-o",
				fmt.Sprintf(
					"minibin_release_%s_%s_%s%s",
					strings.ReplaceAll(version, ".", "-"),
					output.Platform, output.Arch,
					ext,
				),
			)
			cmd.Env = append(os.Environ(),
				"GOOS="+output.Platform,
				"GOARCH="+output.Arch,
			)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
			fmt.Printf(" Done\n")
		}
		cmd := exec.Command("git", "show-ref", "--tags")
		tagsStr, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
		tagFound := false
		tags := strings.Split(string(tagsStr), "\n")
		for i := range tags {
			tag := tags[i]
			if strings.Contains(tag, "tags/"+version) {
				tagFound = true
				break
			}
		}
		if !tagFound {
			fmt.Println("Creating a git tag")
			cmd := exec.Command("git", "tag", version)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Git tag already exists")
		}
		cmd = exec.Command("git", "push", "origin", version)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		cmd = exec.Command(
			"gh", "release", "create", version,
			"minibin_release_*",
			"--title", "\"Release "+version+"\"",
			"--generate-notes",
		)
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Cleanup...")
		for i := range buildOutputs {
			output := buildOutputs[i]
			ext := ""
			if output.Platform == "windows" {
				ext = ".exe"
			}
			filename :=
				fmt.Sprintf(
					"minibin_release_%s_%s_%s%s",
					strings.ReplaceAll(version, ".", "-"),
					output.Platform, output.Arch,
					ext,
				)
			if _, err := os.Stat(filename); err == nil {
				err := os.Remove(filename)
				if err != nil {
					fmt.Println(err)
				}
			} else if os.IsNotExist(err) {
				continue
			} else {
				continue
			}
		}
		fmt.Printf(" Done\n")
		return
	}
	if args.Clear {
		files := fs.ListFiles("")
		if files.Error != nil {
			panic(files.Error)
		}
		for i := range files.Value {
			filePath := files.Value[i]
			if strings.HasPrefix(filePath, "obj_") ||
				strings.HasPrefix(filePath, "enum_") ||
				strings.HasPrefix(filePath, "minibin__") {
				fs.Delete(strings.Split(filePath, fs.Slash)...)
			}
		}
		return
	}
	schemas := schema.Read(args.Input)
	var output []*parser.ParserOutputItem
	switch args.Lang {
	case "go":
		output = parser_go.Parse(schemas, &args)
	case "ts":
		output = parser_ts.Parse(schemas, &args)
	default:
		panic(
			fmt.Errorf("Invalid language provided: %s", args.Lang),
		)
	}
	for i := range output {
		item := output[i]
		fs.Write([]byte(item.Content), item.Path)
	}
}

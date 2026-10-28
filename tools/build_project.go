package tools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bbanez/minibin/src/utils"
)

const Version = "v0.5.0"

type BuildOutputs struct {
	Name     string
	Platform string
	Arch     string
}

func BuildProject(args *utils.Args) {
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
				strings.ReplaceAll(Version, ".", "-"),
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
	if !args.CreateRelease {
		return
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
		if strings.Contains(tag, "tags/"+Version) {
			tagFound = true
			break
		}
	}
	if !tagFound {
		fmt.Println("Creating a git tag")
		cmd := exec.Command("git", "tag", Version)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Git tag already exists")
	}
	cmd = exec.Command("git", "push", "origin", Version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	cmd = exec.Command(
		"gh", "release", "create", Version,
		"minibin_release_*",
		"--title", "\"Release "+Version+"\"",
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
				strings.ReplaceAll(Version, ".", "-"),
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
}

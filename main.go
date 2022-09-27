package main

import (
	"flag"
	"fmt"
	"log"
	"mdj-diff/diff"
	"mdj-diff/printer"
	"mdj-diff/types"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	var cwd, file, source, target, colour string
	var help, json, markdown bool

	flag.StringVar(&cwd, "cd", ".", "working directory")
	flag.StringVar(&file, "f", "pbt.mdj", "the file to diff")
	flag.StringVar(&source, "s", "", "the commit/branch you are working on (optional)")
	flag.StringVar(&colour, "colour", "always", "colour options: always, never, auto (unused)")
	flag.StringVar(&target, "t", "HEAD", "the commit/branch to diff against e.g. master")

	flag.BoolVar(&json, "json", false, "output as json")
	flag.BoolVar(&markdown, "markdown", false, "output as markdown")
	flag.BoolVar(&help, "h", false, "print this help menu")

	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	if target == "" {
		flag.Usage()
		os.Exit(0)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var sourceProject, targetProject *types.Project
	var sourceError, targetError error

	go func() {
		defer wg.Done()
		sourceProject, sourceError = readFileFromGit(cwd, file, source)
	}()
	go func() {
		defer wg.Done()
		targetProject, targetError = readFileFromGit(cwd, file, target)
	}()

	wg.Wait()

	if sourceError != nil {
		log.Fatalf("failed to read from git: %v", sourceError)
	} else if targetError != nil {
		log.Fatalf("failed to read from git: %v", targetError)
	}

	res := diff.DiffTables(sourceProject.GetTableMap(), targetProject.GetTableMap())

	if json {
		printer.PrintJson(res)
	} else if markdown {
		diff.COLOR_ENABLED = false
		printer.PrintMarkdown(sourceProject, targetProject, res)
	} else {
		printer.PrintText(sourceProject, targetProject, res)
	}
}

func readFileFromGit(cwd, path, commit string) (*types.Project, error) {
	if commit == "" {
		return readFile(cwd, path)
	}

	cmd := exec.Command("git", "-C", cwd, "show", fmt.Sprintf("%s:%s", commit, path))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create git show command: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to run git show %s: %v", path, err)
	}

	return types.NewProjectFromJson(stdout)
}

func readFile(cwd, path string) (*types.Project, error) {
	file, err := os.Open(filepath.Join(cwd, path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return types.NewProjectFromJson(file)
}

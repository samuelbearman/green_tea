package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type CommandInfo struct {
	originalCommand string
	args            string
	operatingSystem OperatingSystem
	tempFile        *os.File
	toolBytes       []byte
}

func (command *CommandInfo) ReadEmbeddedContent(pattern string) {
	tool, err := content.ReadFile(command.originalCommand)
	if err != nil {
		fmt.Printf("Error, cannot access %s\n", tool)
		log.Fatal(err.Error())
	}

	command.toolBytes = tool

	tempFile, err := ioutil.TempFile("", pattern)
	if err != nil {
		fmt.Printf("Error cant write to TempFile\n")
		log.Fatal(err.Error())
	}

	command.tempFile = tempFile
}

type OperatingSystem string

const (
	Windows   OperatingSystem = "windows"
	Linux     OperatingSystem = "linux"
	Macintosh OperatingSystem = "mac"
	Unknown   OperatingSystem = "unknown"
)

//go:embed tools/*
var content embed.FS

func main() {

	listCommand := flag.NewFlagSet("list", flag.ContinueOnError)
	runCommand := flag.NewFlagSet("run", flag.ContinueOnError)
	cleanCommand := flag.NewFlagSet("clean", flag.ContinueOnError)

	programToRunPtr := runCommand.String("program", "", "Program to run (Required)")
	argsToRunPtr := runCommand.String("args", "", "Args to run with (Optional)")

	if len(os.Args) < 2 {
		fmt.Println("pass 'list' or 'run' as subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCommand.Parse([]string{})
	case "run":
		runCommand.Parse(os.Args[2:])
	case "clean":
		cleanCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if cleanCommand.Parsed() {
		cleanUp()
	}

	if listCommand.Parsed() {
		listTools()
	}

	if runCommand.Parsed() {
		if *programToRunPtr == "" {
			fmt.Println("Must pass a valid program")
			os.Exit(0)
		} else {
			os := getOperatingSystem()
			commandInfo := CommandInfo{originalCommand: *programToRunPtr, args: *argsToRunPtr, operatingSystem: os, tempFile: nil}
			runTool(&commandInfo)
		}
	}
}

func cleanUp() {
	currentOs := getOperatingSystem()
	if currentOs == "linux" {
		files, err := filepath.Glob("/tmp/gt-*")
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if err := os.Remove(f); err != nil {
				panic(err)
			}
		}
	}
}

func getOperatingSystem() OperatingSystem {
	os := runtime.GOOS

	switch os {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "darwin":
		return "mac"
	default:
		return "unknown"
	}
}

func runTool(command *CommandInfo) {
	pattern := "gt-*"

	switch filepath.Ext(command.originalCommand) {
	case ".sh":
		if command.operatingSystem != Linux {
			log.Fatal("Can only be run on a Linux machine")
		} else {
			pattern += ".sh"
		}
	case ".exe":
		if command.operatingSystem != Windows {
			log.Fatal("Can only be run on a Windows machine")
		} else {
			pattern += ".exe"
		}
	case ".py":
		log.Fatal("Not yet implemented")
	default:
		log.Fatalf("Unable to determine extension for %s", filepath.Ext(command.originalCommand))
	}

	command.ReadEmbeddedContent(pattern)

	os.Chmod(command.tempFile.Name(), 0o500)
	command.tempFile.Write(command.toolBytes)
	command.tempFile.Close()

	runCommand(command)

}

func runCommand(command *CommandInfo) {
	var cmd *exec.Cmd

	switch command.operatingSystem {
	case Windows:
		cmd = exec.Command("cmd", command.tempFile.Name())
	case Linux:
		cmd = exec.Command("/bin/sh", command.tempFile.Name())
	case Macintosh:
		cmd = exec.Command("/bin/sh", command.tempFile.Name())
	case Unknown:
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func listTools() {
	result, err := fs.ReadDir(content, "tools")
	if err != nil {
		log.Fatal(err)
	}

	for _, element := range result {
		if !element.IsDir() {
			fmt.Printf("tools/%s\n", element.Name())
		}
	}
}

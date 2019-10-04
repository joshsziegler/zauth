// +build ignore

package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// runCommand runs the command 'name' with the provided arguments
func runCommand(name string, arg ...string) {
	outBytes, err := exec.Command(name, arg...).CombinedOutput()
	outStr := strings.TrimSpace(string(outBytes))
	if len(outStr) > 0 {
		fmt.Printf("%s\n", outStr)
	}
	if err != nil { // Fail and exit the program now if there was an error
		log.Fatalln(err)
	}
}

func getGitVersion() string {
	outBytes, err := exec.Command("git", "describe", "--dirty").CombinedOutput()
	if err != nil { // Fail and exit the program now if there was an error
		log.Fatalln(err)
	}
	version := strings.TrimSpace(string(outBytes))
	return version
}

func main() {
	fmt.Println("1. Format source files")
	fmt.Println("  - Go")
	runCommand("goimports", "-w", ".")
	// TODO: Format CSS and JavaScript?
	fmt.Println("2. Compile source files")
	fmt.Println("  - Go (into an executable)")
	buildDate := time.Now().Format(time.RFC3339)
	buildVers := getGitVersion()
	runCommand("packr", "build", "-ldflags", "-X main.Version="+buildVers+" -X main.BuildDate="+buildDate)
	fmt.Println("\nBUILD COMPLETE")
}

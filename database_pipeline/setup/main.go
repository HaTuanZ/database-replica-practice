package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

func main() {
	// Execute the mysql command
	cmd := exec.Command("docker", "exec", "-i", "db-master-1", "mysql", "-uroot", "-p9fUcd2^=;V]M", "-e", "show binary log status\\G;")
	// Capture stdout and stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		fmt.Println("Command stderr:", stderr.String())
		return
	}

	// Capture the output
	output := out.String()

	// Use regex to extract the File and Position values
	fileRegex := regexp.MustCompile(`File:\s+(.*)`)
	positionRegex := regexp.MustCompile(`Position:\s+(\d+)`)

	fileMatch := fileRegex.FindStringSubmatch(output)
	positionMatch := positionRegex.FindStringSubmatch(output)

	if len(fileMatch) > 1 && len(positionMatch) > 1 {
		file := fileMatch[1]
		position := positionMatch[1]
		fmt.Println("File:", file)
		fmt.Println("Position:", position)
	} else {
		fmt.Println("Could not extract File and Position from output")
	}
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var err error
	if len(os.Args) == 2 {
		err = runFile(os.Args[1])
	} else if len(os.Args) == 1 {
		err = runPrompt()
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s [script]\n", os.Args[0])
		os.Exit(64)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func runFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	return run(file)
}

func runPrompt() error {
	inputReader := bufio.NewReader(os.Stdin)
	reader := strings.NewReader("")
	for {
		fmt.Printf("> ")
		line, err := inputReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil // done
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		reader.Reset(line)
		err = run(reader)
		if err != nil {
			//TODO call handler
		}
	}
}

func run(reader io.Reader) error {
	s := bufio.NewScanner(reader)
	s.Split(bufio.ScanWords)
	//s.Init(reader)
	for s.Scan() {
		fmt.Println(s.Text())
	}
	return nil
}

func failed(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[%d] Error %s: %s", line, where, message)
}

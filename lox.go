package main

import (
	"bufio"
	"fmt"
	"github.com/keatonmartin/golox/scanner"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		err := runFile(args[0])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		runPrompt()
	}
}

func runFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	run(data)
	return nil
}

func runPrompt() {
	s := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for s.Scan() {
		run(s.Bytes())
		fmt.Print("> ")
	}
}

func run(source []byte) {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()

	for i := 0; i < len(tokens); i++ {
		fmt.Println(tokens[i].String())
	}

	if len(s.Errs) >= 1 {
		for i := 0; i < len(s.Errs); i++ {
			err := s.Errs[i]
			newError(err.Line, err.Message)
		}
		return
	}

}

func newError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Printf("[%d] Error %s: %s\n", line, where, message)
}

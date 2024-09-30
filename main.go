package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"nicholasq.xyz/monkey/evaluator"
	"nicholasq.xyz/monkey/lexer"
	"nicholasq.xyz/monkey/object"
	"nicholasq.xyz/monkey/parser"
	"nicholasq.xyz/monkey/repl"
)

const fileExtension = "mky"

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: %s [script]\n", fileExtension)
		os.Exit(64)
	}

	if len(os.Args) == 2 {
		fileName := os.Args[1]
		runFile(fileName)
	} else {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

func runFile(fileName string) {

	if !strings.HasSuffix(fileName, "."+fileExtension) {
		fmt.Printf("File must be a .%s file.\n", fileExtension)
		os.Exit(1)
	}

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}

	defer file.Close()
	bytes, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("Error reading file: ", err)
		os.Exit(1)
	}

	strContents := string(bytes)
	runScript(strContents)
}

func runScript(script string) {
	env := object.NewEnvironment()

	l := lexer.New(script)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
	}

	evaluator.Eval(program, env)
}

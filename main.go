package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"weaver/lexer"
	"weaver/parser"
)

func processFile(filename string) (*parser.FileNode, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s : %v", filename, err)
	}
	l := lexer.NewLexer(bytes.NewReader(data), false)

	p := parser.NewParser(l)

	astFile, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return astFile, nil
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Please provide a .wevr file to parse, weaver <file.wevr>")
	}
	filename := os.Args[1]

	astFile, err := processFile(filename)
	if err != nil {
		log.Fatalf("Error processing file %s: %v", filename, err)
	}

	fmt.Println(filename, "parsed successfully. AST:")
	jsonData, err := json.MarshalIndent(astFile, "", " ")
	if err != nil {
		log.Fatalf("Error marshalling AST to JSON: %v", err)
	}

	fmt.Println(string(jsonData))

}

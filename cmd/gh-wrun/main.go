package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/t4kamura/gh-wrun/internal/input"
)

const version = "0.0.0"

func main() {
	v := flag.Bool("v", false, "show version")
	b := flag.Bool("b", false, "first interactively select a git branch name")
	flag.Parse()

	if *v {
		fmt.Printf("gh-wrun version %s\n", version)
		os.Exit(0)
	}

	r, err := input.NewInputResult(!*b)

	if err != nil {
		log.Fatal(err)
	}

	if !r.IsRun {
		log.Fatal("Canceled")
	}

	if err := r.Workflow.Run(r.Branch, r.WorkflowInputs); err != nil {
		log.Fatal(err)
	}
}
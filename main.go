package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version = "0.0.0"

func main() {
	v := flag.Bool("v", false, "show version")
	flag.Parse()

	if *v {
		fmt.Printf("gh-wrun version %s\n", version)
		os.Exit(0)
	}

	r, err := NewInputResult()

	if err != nil {
		log.Fatal(err)
	}

	if !r.isRun {
		log.Fatal("Canceled")
	}

	if err := r.workflow.Run(*r.branch, *r.workflowInputs); err != nil {
		log.Fatal(err)
	}
}

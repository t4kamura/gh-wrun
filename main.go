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
	l := flag.Bool("l", false, "check workflow file with linter before parsing")
	flag.Parse()

	if *v {
		fmt.Printf("gh-wrun version %s\n", version)
		os.Exit(0)
	}

	r, err := NewInputResult(*l)

	if err != nil {
		log.Fatal(err)
	}

	if !r.isRun {
		log.Fatal("Canceled")
	}

	if err := r.workflow.Run(r); err != nil {
		log.Fatal(err)
	}
}

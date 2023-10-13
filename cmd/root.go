package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/t4kamura/gh-wrun/internal/input"
	ver "github.com/t4kamura/gh-wrun/internal/version"
)

const (
	version           = "0.0.0"
	requiredGhVersion = "2.35.0"
)

func Execute() {
	v := flag.Bool("v", false, "show version")
	b := flag.Bool("b", false, "first interactively select a git branch name")
	flag.Parse()

	if *v {
		fmt.Printf("gh-wrun version %s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) != 0 {
		flag.Usage()
		os.Exit(1)
	}

	result, err := ver.CheckGhVersion(requiredGhVersion)
	if err != nil {
		log.Fatal(err)
	}

	if !result {
		log.Fatalf("gh-wrun requires gh version %s or later", requiredGhVersion)
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

	fmt.Println("Workflow started")
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hongjunxin/go-learning/os/flag/internal"
)

var (
	name string
)

func init() {
	flag.CommandLine = flag.NewFlagSet("question", flag.ExitOnError)
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "question")
		flag.PrintDefaults()
	}
	flag.StringVar(&name, "name", "everyone", "The greeting object.")
}

func main() {
	flag.Parse()
	fmt.Printf("name: %v\n", name)
	internal.Hello()
}

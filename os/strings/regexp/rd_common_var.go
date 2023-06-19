package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	f, err := os.ReadFile("tmp/common-var.conf")
	if err != nil {
		log.Fatalf("read common-var.conf failed, err=%v", err)
	}

	exp, err := regexp.Compile("MY_PROJECT_ENV_NAME=\"(.+)\"")
	if err == nil {
		matchAll := exp.FindAllSubmatch(f, -1)
		fmt.Printf("len(matchAll)=%d\n", len(matchAll))
		for _, m := range matchAll {
			for _, n := range m {
				fmt.Println(string(n))
			}
		}
		fmt.Println()

		match := exp.FindSubmatch(f)
		fmt.Printf("len(match)=%d\n", len(match))
		for _, m := range match {
			fmt.Println(string(m))
		}
	}
}

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	name    string
	version string
	gitSHA  string
)

const usage = `version: %s
Usage: %[2]s [-v] [-h] <PATH>
Options:
    -h        help
    -v        show version and exit
Examples:
    %[2]s .
    %[2]s /home/user/go/src/github.com/username/repo
`

var erroneousWords = [][]byte{
	[]byte("golang"),
	[]byte("Golang"),
	[]byte("goLang"),
	[]byte("GoLang"),
	[]byte("GOLANG"),
}

const properName = "Go"

func main() {
	var vers bool

	flag.Usage = func() {
		w := os.Stderr
		for _, arg := range os.Args {
			if arg == "-h" {
				w = os.Stdout
				break
			}
		}
		fmt.Fprintf(w, usage, version, name)
	}

	flag.BoolVar(&vers, "v", false, "")
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stdout, "version: %s - git sha: %s\n", version, gitSHA)
		return
	}

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, usage, version, name)
		os.Exit(1)
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".md") {
			fmt.Printf("found: %s\n", path)

			b, err := os.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, ew := range erroneousWords {
				b = bytes.Replace(b, ew, []byte(properName), -1)
			}

			if err = os.WriteFile(path, b, info.Mode()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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
Usage: %[2]s [-v] [-h] [-d] <PATH>
Options:
    -h        help
    -v        show version and exit
    -d        dry-run prints file names with found words
Examples:
    %[2]s .
    %[2]s /home/user/go/src/github.com/username/repo
	%[2]s -d /home/user/go/src/github.com/username/repo
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
	var dryRun string

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
	flag.StringVar(&dryRun, "d", "", "")
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stdout, "version: %s - git sha: %s\n", version, gitSHA)
		return
	}

	var path string
	if dryRun != "" {
		path = dryRun
	} else {
		path = os.Args[1]
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(p, ".md") {
			b, err := os.ReadFile(p)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, ew := range erroneousWords {
				if dryRun != "" {
					if bytes.Contains(b, ew) {
						fmt.Printf("%s in %s\n", string(ew), p)
					}
					continue
				}

				b = bytes.Replace(b, ew, []byte(properName), -1)
				if err = os.WriteFile(p, b, info.Mode()); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if err = os.WriteFile(p, b, info.Mode()); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//go:build go1.21

/*
MIT License

Copyright (c) 2026 Olivier Mengué and contributors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Command readme-gofmt applies 'gofmt -s' to all Go block in README.md.
package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"regexp"
	"slices"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("readme-gofmt: ")

	buf, err := os.ReadFile("README.md")
	if err != nil {
		log.Fatal(err)
	}

	gofmt, err := exec.LookPath("gofmt")
	if err != nil {
		log.Fatal(err)
	}

	buf = bytes.ReplaceAll(buf, []byte("\r"), nil) // CRLF -> LF

	reBlock := regexp.MustCompile("(?s)\n```go\n(.*?\n)```")

	matches := reBlock.FindAllSubmatchIndex(buf, -1)

	changes := 0

	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		block := buf[match[2]:match[3]]

		cmd := exec.Command(gofmt, "-s")
		cmd.Stdin = bytes.NewReader(block)
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		outputBytes := bytes.ReplaceAll(output.Bytes(), []byte("\r"), []byte("")) // CRLF -> LF
		if !bytes.Equal(outputBytes, block) {
			changes++
			buf = slices.Replace(buf, match[2], match[3], outputBytes...)
		}
	}

	if changes > 0 {
		cmd := exec.Command("diff", "-a", "-u", "README.md", "-")
		cmd.Stdin = bytes.NewReader(buf)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		os.Exit(2)
	}
}

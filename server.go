package main

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
)

func responseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "File %q\n\n", html.EscapeString(r.URL.Path[1:]))

	fmt.Fprintf(w, "File contents:\n")

	fileToRead := r.URL.Path[1:]

	fmt.Fprintf(w, "beg file\n===========================\n")

	buffer := make([]byte, 0)
	size := readFile(fileToRead, buffer)

	if _, err := w.Write(buffer[:size]); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "\n===========================\nend")
}

func readFile(filename string, buffer []byte) int {
	if filename != "favicon.ico" {
		fileIn, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := fileIn.Close(); err != nil {
				panic(err)
			}
		}()
		reader := bufio.NewReader(fileIn)

		var size int = 0
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}

			if n == 0 {
				break
			}
			size += n
			copy(buffer, buf)
		}
		return size
	}
	return 0
}

func main() {
	http.HandleFunc("/", responseHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

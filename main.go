package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"strings"
)

func emailSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// A single blank line and a "From " separates a message
	// https://en.wikipedia.org/wiki/Mbox#Family
	if i := strings.Index(string(data), "\n\nFrom "); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return
}

func readEmail(b []byte) {
	// To properly read a mail message, we need to remove any preceeding
	// newlines and additionally remove the "From " line
	const NL = "\n"
	trimmed := strings.TrimLeft(string(b), NL)
	msgString := strings.Join(strings.Split(trimmed, NL)[1:], NL)

	msg, err := mail.ReadMessage(strings.NewReader(msgString))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("From:", msg.Header.Get("From"))
}

func emailScanner(mbox io.Reader) {
	// An artificial input source.
	scanner := bufio.NewScanner(mbox)

	// Allow a maximum of 2^24 bytes per message
	scanner.Buffer([]byte{}, 1<<24)
	scanner.Split(emailSplit)

	count := 0
	for scanner.Scan() {
		count++
		readEmail(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	fmt.Println("Total emails:", count)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage:", os.Args[0], "<filename>")
	}

	filename := os.Args[1]
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		log.Fatalln("Unable to open file:", err)
	}

	emailScanner(f)

}

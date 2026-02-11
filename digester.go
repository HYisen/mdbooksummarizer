package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
)

var Digester DigesterFunc = Header

type DigesterFunc func(content []byte) (title string)

// FileSize use content bytes length as title for debug purpose.
//
//goland:noinspection GoUnusedExportedFunction
func FileSize(content []byte) (title string) {
	return fmt.Sprintf("filesize %d", len(content))
}

func Header(content []byte) (title string) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Scan()
	return ParseMdHeader(scanner.Text())
}

func ParseMdHeader(s string) string {
	if !strings.HasPrefix(s, "# ") {
		log.Fatalf("unexpected header line [%s]", s)
	}
	return strings.TrimLeft(s, "# ")
}

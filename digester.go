package main

import "fmt"

var Digester = FileSize

type DigesterFunc func(content []byte) (title string)

func FileSize(content []byte) (title string) {
	return fmt.Sprintf("filesize %d", len(content))
}

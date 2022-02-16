package main

import (
	"os"

	"github.com/adelgado0723/processor"
)

func main() {
	pipeline := processor.NewPipeline(os.Stdin, os.Stdout, nil, 64)

}

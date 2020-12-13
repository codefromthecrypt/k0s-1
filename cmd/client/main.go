package main

import (
	"log"
	"os"

	"k0s.io/k0s/pkg/cli/client"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Fatalln(client.Run(os.Args[1:]))
}

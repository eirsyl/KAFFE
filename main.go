package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/webkom/KAFFE/cmd"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

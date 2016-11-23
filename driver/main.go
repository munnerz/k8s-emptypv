package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"io"
)

type FlexStatus string

const (
	FlexStatusSuccess FlexStatus = "Success"
	FlexStatusFailure FlexStatus = "Failure"
)

var (
	cmds = map[string]CmdFn{
		"init": initDriver,
		"attach": attach,
		"mount": mount,
		"detach": detach,
		"unmount": unmount,
	}
)

type CmdFn func([]string) (FlexOutput, error)

type FlexOutput struct {
	Status  FlexStatus `json:"status"`
	Message string     `json:"message"`
	Device  string     `json:"device,omitempty"`
}

type FlexOptions struct {
	VolumeID    string `json:"volumeID"`
	Path        string `json:"path"`
	Permissions os.FileMode `json:"permissions,omitempty"`
	ReadWrite   string `json:"readWrite"`
	FSType      string `json:"fsType"`
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("invalid command input: %s", os.Args)
	}

	cmd := os.Args[1]
	var fn CmdFn
	var ok bool
	if fn, ok = cmds[cmd]; !ok {
		log.Fatalf("invalid command: %s", cmd)
	}

	out, err := fn(os.Args)

	if err != nil {
		defer os.Exit(1)
	}

	// we always write output
	err = writeOutput(out, os.Stdout)

	if err != nil {
		fmt.Printf("error writing output to stdout: %s", err.Error())
		os.Exit(2)
	}
}

func writeOutput(o FlexOutput, out io.Writer) error {
	enc := json.NewEncoder(out)
	return enc.Encode(o)
}
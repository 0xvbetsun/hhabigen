package main

import (
	"fmt"
	"os/exec"
)

type ABI struct {
	ContractName     string      `json:"contractName"`
	SourceName       string      `json:"sourceName"`
	ABI              interface{} `json:"abi"`
	ByteCode         string      `json:"bytecode"`
	DeployedByteCode string      `json:"deployedBytecode"`
}

func run(in, typ, pkg, output string) error {
	args := []string{
		"-abi", in,
		"-out", output,
		"-type", typ,
		"-pkg", pkg,
	}

	cmd := exec.Command("abigen", args...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("generate bindings: %w", err)
	}

	return nil
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func extract(in io.Reader) (*ABI, error) {
	bytes, err := io.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("read artifact: %w", err)
	}

	abi := &ABI{}
	if err = json.Unmarshal(bytes, abi); err != nil {
		return nil, fmt.Errorf("unmarshal artifact: %w", err)
	}

	if abi.ABI == nil || len(abi.ABI.([]interface{})) == 0 || abi.ByteCode == "" || abi.ByteCode == "0x" {
		return nil, errors.New("invalid artifact")
	}

	return abi, nil
}

func extractFromFile(in, out string) (string, error) {
	path, err := filepath.Abs(in)
	if err != nil {
		return "", nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer fp.Close()

	abi, err := extract(fp)
	if err != nil {
		return "", fmt.Errorf("extract abi from %s: %w", path, err)
	}

	if abi == nil {
		return "", errors.New("abi is empty")
	}

	fpOut, err := os.Create(out)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer fpOut.Close()

	marshalled, err := json.Marshal(abi.ABI)
	if err != nil {
		return "", fmt.Errorf("marshal abi: %w", err)
	}

	if _, err = fpOut.Write(marshalled); err != nil {
		return "", fmt.Errorf("write abi to disk: %w", err)
	}

	return abi.ContractName, nil
}

func processFile(inFile, outDir string) error {
	ext := filepath.Ext(inFile)
	if ext != ".json" {
		return fmt.Errorf("supported only json files, got %s", ext)
	}

	abiOutDir, err := mkDirIfNotExist(outDir, AbiDir)
	if err != nil {
		return fmt.Errorf("make out dir: %w", err)
	}

	abiOutFile := path.Join(abiOutDir, filepath.Base(inFile))

	contractName, err := extractFromFile(inFile, abiOutFile)
	if err != nil {
		return fmt.Errorf("extract abi: %w", err)
	}

	if len(contractName) > 0 {
		bindingsOutDir, err := mkDirIfNotExist(outDir, *out)
		if err != nil {
			return fmt.Errorf("make out dir: %w", err)
		}
		bindingsOutFile := path.Join(bindingsOutDir, fmt.Sprintf("%s.go", contractName))
		fmt.Printf("generating bindings for %s", contractName)

		if err = run(abiOutFile, contractName, *pkg, bindingsOutFile); err != nil {
			return fmt.Errorf("abigen: %w", err)
		}
	}

	return nil
}

func process(f, o string, isDir bool) error {
	if !isDir {
		return processFile(f, o)
	}

	return filepath.Walk(f, func(p string, info fs.FileInfo, err error) error {
		if filepath.Ext(p) == ".json" {
			return processFile(p, o)
		}
		return nil
	})
}

func mkDirIfNotExist(prefix, dir string) (string, error) {
	newPath := path.Join(prefix, dir)
	info, err := os.Stat(newPath)

	if err != nil {
		err = os.MkdirAll(newPath, 0o755)
		if err != nil {
			return "", fmt.Errorf("create dir: %w", err)
		}
	} else {
		if !info.IsDir() {
			return "", errors.New("output path already exists, but is not a directory.")
		}
	}

	return newPath, nil
}

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	algorithmName := flag.String("name", "", "The algorithm's name")
	algorithmDesc := flag.String("desc", "", "The algorithm's description")
	algorithmAuthor := flag.String("author", "", "The algorithm's author")
	algorithmLicense := flag.String("license", "", "The algorithm's license")
	remoteURL := flag.String("website", "", "The algorithm's website")
	downloadUrl := flag.String("download-url", "", "The algorithm's website")
	initFunctions := flag.String("init", "", "comma deliminated list of initialization funcs in order (for example -init-func _start,initialize)")
	allocatingFunction := flag.String("allocating", "", "The memory allocting func (for example, malloc)")
	freeingFunction := flag.String("freeing", "", "The memory freeing func (for example, free)")
	wasm := flag.String("wasm", "", "The wasm file name")
	moduleName := flag.String("mod", "", "The module name (for example, env)")
	version := flag.Int("version", 0, "The version of the plugin")
	apiVersion := flag.Int("api", 1, "The version of the plugin api")
	outfile := flag.String("o", "", "The outfile")
	flag.Parse()
	if *algorithmName == "" {
		log.Fatal("The algorithm name cannot be empty")
	}
	if *algorithmAuthor == "" {
		log.Fatal("The algorithm author cannot be empty")
	}
	if *algorithmLicense == "" {
		log.Fatal("The algorithm license cannot be empty")
	}
	if *remoteURL == "" {
		log.Fatal("The remote website cannot be empty")
	}
	if *initFunctions == "" {
		log.Fatal("The init functions list cannot be empty")
	}
	if *allocatingFunction == "" {
		log.Fatal("The allocating functions cannot be empty")
	}
	if *freeingFunction == "" {
		log.Fatal("The freeing functions cannot be empty")
	}
	if *wasm == "" {
		log.Fatal("The wasm file cannot be empty")
	}
	if *moduleName == "" {
		log.Fatal("The module name cannot be empty")
	}
	if *version == 0 {
		log.Fatal("The version cannot be 0")
	}
	if *outfile == "" {
		log.Fatal("The output file cannot be empty")
	}

	fmt.Println(*wasm)
	conts, err := os.ReadFile(*wasm)
	if err != nil {
		panic(err)
	}
	wasmEncoded := base64.StdEncoding.EncodeToString(conts)

	data, err := json.Marshal(map[string]any{
		"name":         *algorithmName,
		"desc":         *algorithmDesc,
		"author":       *algorithmAuthor,
		"license":      *algorithmLicense,
		"remote-url":   *remoteURL,
		"download-url": *downloadUrl,
		"init":         *initFunctions,
		"alloc":        *allocatingFunction,
		"dealloc":      *freeingFunction,
		"wasm":         wasmEncoded,
		"module-name":   *moduleName,
		"version":      *version,
		"api-version":  *apiVersion,
	})
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(*outfile, []byte((data)), 0755)
	if err != nil {
		panic(err)
	}
}

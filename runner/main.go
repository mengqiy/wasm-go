package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	wasmtime "github.com/bytecodealliance/wasmtime-go"

	// Not using github.com/wasmerio/wasmer-go, since it doesn't provide a convenient way to configure stdin and stdout.
)

func runOneWASI(wasmFile, stdinPath, stdoutPath string) {
	// TODO: we can reuse the same WASM engine.
	engine := wasmtime.NewEngine()

	// Create our module
	wasm, err := ioutil.ReadFile(wasmFile)
	check(err)
	module, err := wasmtime.NewModule(engine, wasm)
	check(err)

	// Create a linker with WASI functions defined within it
	linker := wasmtime.NewLinker(engine)
	err = linker.DefineWasi()
	check(err)

	// Configure WASI imports to write stdout into a file, and then create
	// a `Store` using this wasi configuration.
	wasiConfig := wasmtime.NewWasiConfig()
	// Use file as stdin and write stdout to a different file.
	wasiConfig.SetStdinFile(stdinPath)
	wasiConfig.SetStdoutFile(stdoutPath)

	store := wasmtime.NewStore(engine)
	store.SetWasi(wasiConfig)
	instance, err := linker.Instantiate(store, module)
	check(err)

	// Run the function
	nom := instance.GetFunc(store, "_start")
	_, err = nom.Call(store)
	check(err)
}

func main() {
	dir, err := ioutil.TempDir("", "wasi-out")
	check(err)
	defer os.RemoveAll(dir)

	stdoutPath := filepath.Join(dir, "file2")
	input, err := ioutil.ReadAll(os.Stdin)
	check(err)
	stdinPath := filepath.Join(dir, "file1")
	ioutil.WriteFile(stdinPath, input, 0644)

	// The following 2 functions work with yaml.
	// rustSetNsWasm := "../fn-set-ns-rust/target/wasm32-wasi/release/set_ns.wasm"
	// rustNamePrefixWasm := "../fn-name-prefix-rust/target/wasm32-wasi/release/name_prefix.wasm"
	// wasmList := []string{rustSetNsWasm, rustNamePrefixWasm}

	// The following function work with json.
	goSetNsWasm := "../fn-set-ns-go/fn-set-ns-go.wasm"
	wasmList := []string{goSetNsWasm}

	for _, wasm := range wasmList {
		runOneWASI(wasm, stdinPath, stdoutPath)
		stdinPath, stdoutPath = stdoutPath, stdinPath
	}

	stdinPath, stdoutPath = stdoutPath, stdinPath
	out, err := ioutil.ReadFile(stdoutPath)
	check(err)
	fmt.Println(string(out))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

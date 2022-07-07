package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

func runOneWASI(wasmFile, stdinPath, stdoutPath string) {
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

	setNsWasm := "/usr/local/google/home/mengqiy/rustprojects/set_ns/target/wasm32-wasi/release/set_ns.wasm"
	namePrefixWasm := "/usr/local/google/home/mengqiy/rustprojects/name_prefix/target/wasm32-wasi/release/name_prefix.wasm"
	wasmList := []string{setNsWasm, namePrefixWasm}
	for _, wasm := range wasmList {
		runOneWASI(wasm, stdinPath, stdoutPath)
		stdinPath, stdoutPath = stdoutPath, stdinPath
	}

	stdinPath, stdoutPath = stdoutPath, stdinPath
	out, err := ioutil.ReadFile(stdoutPath)
	check(err)
	fmt.Println(string(out))


	// input, err := ioutil.ReadAll(os.Stdin)
	// check(err)

	// // setNsWasmBytes, _ := ioutil.ReadFile("/usr/local/google/home/mengqiy/rustprojects/set_ns/target/wasm32-wasi/release/set_ns.wasm")
	// setNsWasmBytes, _ := ioutil.ReadFile("/usr/local/google/home/mengqiy/rustprojects/set_ns/target/wasm32-wasi/debug/set_ns.wasm")
	// // namePrefixWasmBytes, _ := ioutil.ReadFile("/usr/local/google/home/mengqiy/rustprojects/name_prefix/target/wasm32-wasi/release/name_prefix.wasm")

	// store := wasmer.NewStore(wasmer.NewEngine())
	// setNsModule, _ := wasmer.NewModule(store, setNsWasmBytes)
	// // namePrefixModule, _ := wasmer.NewModule(store, namePrefixWasmBytes)

	// bytes, err := setNsModule.Serialize()
	// fmt.Printf("setNsModule: %v, err: %v\n", string(bytes), err)

	// wasiEnv, _ := wasmer.NewWasiStateBuilder("wasi-program").
	// 	Finalize()
	// setNsImportObject, err := wasiEnv.GenerateImportObject(store, setNsModule)
	// check(err)
	// // namePrefixImportObject, err := wasiEnv.GenerateImportObject(store, namePrefixModule)
	// // check(err)

	// instance, err := wasmer.NewInstance(setNsModule, setNsImportObject)
	// check(err)

	// exportsMap := instance.Exports.GetExports()
	// for name, extern := range(exportsMap) {
	// 	fmt.Printf("name: %v: extern type: %v, kind: %v\n", name, extern.Type(), extern.Kind())
	// }
	// hb, err := instance.Exports.GetGlobal("__heap_base")
	// val, err := hb.Get()
	// fmt.Printf("global: %v, err: %v\n", val, err)
	// fmt.Printf("global type: %#v\n", hb.Type())

	// start, err := instance.Exports.GetWasiStartFunction()
	// check(err)
	// start()

	// setNsTransformer, err := instance.Exports.GetFunction("set_ns_transformer")
	// check(err)
	// intermediate, err := setNsTransformer(string(input))
	// check(err)

	// // namePrefixTransformer, err := instance.Exports.GetFunction("name_prefix_transformer")
	// // check(err)
	// // result, err := namePrefixTransformer(intermediate)
	// // check(err)

	// fmt.Println(intermediate)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// func main() {
// 	wasmBytes, _ := ioutil.ReadFile("/usr/local/google/home/mengqiy/rustprojects/set_ns/target/wasm32-wasi/release/set_ns.wasm")

// 	store := wasmer.NewStore(wasmer.NewEngine())
// 	module, _ := wasmer.NewModule(store, wasmBytes)

// 	wasiEnv, _ := wasmer.NewWasiStateBuilder("wasi-program").
// 		InheritStdin().
// 		InheritStdout().
// 		InheritStderr().
// 		Finalize()
// 	importObject, err := wasiEnv.GenerateImportObject(store, module)
// 	check(err)

// 	instance, err := wasmer.NewInstance(module, importObject)
// 	check(err)

// 	start, err := instance.Exports.GetWasiStartFunction()
// 	check(err)
// 	start()
// }
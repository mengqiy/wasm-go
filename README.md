# KPT function WASM PoC

In this PoC, we build a simple WASM runner in golang and run multiple KRM functions that are compiled to WASM format.

## Background

WebAssembly is only allowed to be run in a JS host environment. Fortunately, WebAssembly community also provide something
called WebAssembly System Interface (WASI) which provides a way to run WebAssembly on non-JS host environments.
There is an ongoing efforts to standardize it.

## WASM KRM function in golang

Current status in golang:
- golang doesn't support WASI natively and there is no plan to support it: https://github.com/golang/go/issues/31105
- [TinyGo](https://github.com/tinygo-org/tinygo) provides some support for compiling golang program to WASI. But there
  are some significant limitations: the `reflect` package is not fully supported and several pacakges (e.g. `encoding/json`
  and `github.com/go-yaml/yaml`) don't work due to this reason.
  See [issue](https://github.com/tinygo-org/tinygo/issues/2950) and [issue](https://github.com/tinygo-org/tinygo/issues/447)
  for more details.

In `fn-set-ns-go`, I use `github.com/valyala/fastjson` to build a simple KRM function to mutate `metadata.namespace`.

To compile it, you need to ensure you have [`tinygo`](https://github.com/tinygo-org/tinygo) installed.

```shell
$ cd fn-set-ns-go
```

Compile it to WASM

```shell
$ tinygo build -o fn-set-ns-go.wasm -target wasi ./...
```

## WASM KRM function in rust

Rust appears to have better support for WASM than golang.

Please ensure you have `rust`, `cargo` and [`cargo wasi`](https://github.com/bytecodealliance/cargo-wasi) installed.

```shell
$ cd fn-name-prefix-rust
```

Compile it to WASI

```shell
$ cargo wasi build --release
```

The WASI artifact can be found in `target/wasm32-wasi/release/`.

Similiarly for the KRM function under `fn-set-ns-rust` directory.

## WASM KRM function in python

There will be alpha support in python 3.11 whihch is scheduled to release in October 2022.

When we can convert a python program to WASI, it can be a game changer for us.

## WASM runner in golang

### Running with json input

If you are running with the WASM KRM golang functions above, please ensure the following code are uncommented:

```go
goSetNsWasm := "../fn-set-ns-go/fn-set-ns-go.wasm"
wasmList := []string{goSetNsWasm}
```

Execute the runner and pass it a resource list in json format:

```shell
$ cat resourceList.json | go run .
```

You should be able to see the resourceList has been updated.

```json
{"apiVersion":"config.kubernetes.io/v1alpha1","kind":"ResourceList","items":[{"apiVersion":"v1","kind":"ReplicationController","metadata":{"name":"bob","namespace":"prod"},"spec":{"replicas":1,"selector":{"app":"nginx"},"templates":{"metadata":{"name":"nginx","labels":{"app":"nginx"}},"spec":{"containers":[{"name":"nginx","image":"nginx","ports":[{"containerPort":80}]}]}}}},{"apiVersion":"example.com/v1","kind":"MyFoo","metadata":{"name":"bob","namespace":"prod"},"spec":{}}],"functionConfig":{"apiVersion":"v1","kind":"ConfigMap","metadata":{"name":"fn-cfg"},"data":{"namespace":"prod"}}}
```

### Running with yaml input

If you are running with the WASM KRM rust functions above, please ensure the following code are uncommented:

```go
rustSetNsWasm := "../fn-set-ns-rust/target/wasm32-wasi/release/set_ns.wasm"
rustNamePrefixWasm := "../fn-name-prefix-rust/target/wasm32-wasi/release/name_prefix.wasm"
wasmList := []string{rustSetNsWasm, rustNamePrefixWasm}
```

Execute the runner and pass it a resource list in yaml format:

```shell
$ cat resourceList.yaml | go run .
```

You should be able to see the resourceList has been updated.

# Go WASM Test

## Initial Setup

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

## Building

```shell
GOOS=js GOARCH=wasm go build -o main.wasm
```

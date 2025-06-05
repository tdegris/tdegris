package generate

//go:generate bash -c "GOOS=js GOARCH=wasm go build -o ../../../res/main.wasm ../wasm/wasm.go"
//go:generate gzip -f -9 ../../../res/main.wasm

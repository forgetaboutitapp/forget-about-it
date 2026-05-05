module github.com/forgetaboutitapp/forget-about-it/server

go 1.26.2

require (
	github.com/adrg/xdg v0.5.3
	github.com/google/uuid v1.6.0
	github.com/hashicorp/go-set v0.1.14
	github.com/rs/cors v1.11.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/tetratelabs/wazero v1.9.0
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/net v0.47.0
	google.golang.org/protobuf v1.36.11
	modernc.org/sqlite v1.37.0
)

require (
	connectrpc.com/connect v1.19.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	modernc.org/libc v1.62.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.9.1 // indirect
)

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	google.golang.org/protobuf/cmd/protoc-gen-go
)

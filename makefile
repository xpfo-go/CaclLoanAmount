arch = amd64 # amd64 arm64
goos = windows # windows linux darwin

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o caclLoanAmount.exe main.go
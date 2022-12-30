GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/main cmd/lambda/main.go
cd build
zip main.zip main
cd ../cdk
cdk synth && cdk deploy
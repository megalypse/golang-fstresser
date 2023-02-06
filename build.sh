for paths in ./cmd/*/main.go; do
  [[ $paths =~ ./cmd/(.*)/main.go ]]
  binName=${BASH_REMATCH[1]}

  GOARCH=amd64 GOOS=darwin go build -o ./bin/${binName}-darwin ./cmd/${binName}/main.go
  GOARCH=amd64 GOOS=linux go build -o ./bin/${binName}-linux ./cmd/${binName}/main.go
  GOARCH=amd64 GOOS=windows go build -o ./bin/${binName}-windows ./cmd/${binName}/main.go
done

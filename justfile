build:
    go build -ldflags "-s -w" -o build/alfred-picular .
    upx build/alfred-picular

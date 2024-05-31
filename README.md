# Disclaimer

**Before building the application**, you must change the `GOOS` environment variable to match your operating system.

## Build Instructions

### For Linux
```sh
CGO_ENABLED=1 GOOS="linux" GOARCH="amd64" go build -o ./bin/main
```

### For Macos
```sh
@CGO_ENABLED=1 GOOS="darwin" GOARCH="amd64" go build -o ./bin/main```
```

### For Windows
```sh
@CGO_ENABLED=1 GOOS="windows" GOARCH="amd64" go build -o ./bin/main
```

# Usage
#### The app is straightforward to use. Just export the binary to your PATH to use it once you open the terminal.
# Important

#### Make sure to update the .env file to specify the path for the SQLite3 database file location.

# Contribution
#### Any Pull Request to fix things or to add more features would be highly appreciated.

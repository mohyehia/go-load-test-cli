# go-load-test-cli

Small Go CLI for running concurrent load tests against HTTP APIs.

## What it supports

- count-based load profile via `-n`
- duration-based load profile via `-d`
- configurable concurrency via `-c`
- request timeout via `-t`
- custom headers via `-H`
- request body via `-b` or body file via `-f`
- live progress output and final summary

## Build

### Windows

```powershell
go build -o goku.exe .
```

### Linux

```bash
go build -o goku .
```

### macOS

```bash
go build -o goku .
```

## Run

### GET request

#### Windows

```powershell
.\goku.exe -u https://httpbin.org/get -m GET -n 100 -c 10
```

#### Linux

```bash
./goku -u https://httpbin.org/get -m GET -n 100 -c 10
```

#### macOS

```bash
./goku -u https://httpbin.org/get -m GET -n 100 -c 10
```

### POST request with inline JSON body

#### Windows

```powershell
.\goku.exe -u https://httpbin.org/post -m POST -n 100 -c 10 -H "Content-Type: application/json" -b '{"username":"michael"}'
```

#### Linux

```bash
./goku -u https://httpbin.org/post -m POST -n 100 -c 10 -H "Content-Type: application/json" -b '{"username":"michael"}'
```

#### macOS

```bash
./goku -u https://httpbin.org/post -m POST -n 100 -c 10 -H "Content-Type: application/json" -b '{"username":"michael"}'
```

### POST request with body from file

#### Windows

```powershell
.\goku.exe -u https://httpbin.org/post -m POST -n 100 -c 10 -H "Content-Type: application/json" -f .\payload.json
```

#### Linux

```bash
./goku -u https://httpbin.org/post -m POST -n 100 -c 10 -H "Content-Type: application/json" -f ./payload.json
```

#### macOS

```bash
./goku -u https://httpbin.org/post -m POST -n 100 -c 10 -H "Content-Type: application/json" -f ./payload.json
```

## Main flags

- `-u`, `--url` target URL
- `-m`, `--method` HTTP method
- `-n`, `--requests` total requests
- `-d`, `--duration` run duration
- `-c`, `--concurrency` number of workers
- `-t`, `--timeout` per-request timeout
- `-H`, `--header` repeatable HTTP header
- `-b`, `--body` inline request body
- `-f`, `--file` request body file path

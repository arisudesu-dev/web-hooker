web-hooker
==========

A simple web server that runs CGI scripts from working directory. I use it as a server for github webhooks.

## Config

Configuration is done with environment variables.

| Variable | Default | Description |
| -------- | ------- | ----------- |
| PORT | 8000 | TCP port to listen |

## Usage example

Start server at the directory with examples:
```shell
go build -o web-hooker
pushd examples && ../web-hooker
```

Request to run examples/test.sh as CGI script.

```shell
curl -X POST http://localhost:8000/test.sh
```

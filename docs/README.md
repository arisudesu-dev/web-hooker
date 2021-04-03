web-hooker
==========

A simple web server that runs CGI scripts from working directory. I use it as a server for github webhooks.

## Config

Configuration is done with environment variables.

| Variable | Default | Description |
| -------- | ------- | ----------- |
| PORT | 8000 | TCP port to listen |

## Usage example

Start server at directory with scripts:
```shell
go build -o web-hooker
pushd scripts && ../web-hooker
```

Request to run $CWD/execute_hook.sh as CGI script.

```shell
curl -X POST http://localhost:8000/execute_hook.sh
```

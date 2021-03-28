web-hooker
==========

A simple web server that runs CGI scripts from a directory. I use it as a server for github webhooks.

## Config

Configuration is done with environment variables.

| Variable | Default | Description |
| -------- | ------- | ----------- |
| PORT | 8000 | TCP port to listen |
| SCRIPTS_DIR | / | Directory containing launched scripts |

## Usage example

Request to run $SCRIPTS_DIR/execute_hook.sh as CGI script.

```shell
curl -X POST http://localhost:8000/execute_hook.sh
```

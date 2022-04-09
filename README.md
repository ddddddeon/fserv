# fserv
A CLI tool to serve local files via HTTP

Inspired by python's `http.server`/`SimpleHTTPServer` modules, this program serves a directory or file via HTTP. 

## Installation

```bash
git clone https://github.com/ddddddeon/fserv
cd fserv
make
sudo make install
```
## Usage
The server listens on all interfaces, and serves the current working directory by default, on port 8080.

A unique token is generated and output on startup that is required as a query parameter for anyone consuming the endpoint.

```bash
# Serves ./ on http://0.0.0.0:8080?t=<token>
fserv

#Serves /my/directory on http://0.0.0.0:9090?t=<token>
fserv /my/directory 9090
```

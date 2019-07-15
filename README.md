# ft

A file transferer with gRPC.
Inspired by [mattn/ft](https://github.com/mattn/ft).

## Usage

### Run as Server

``` 
$ ft -s
```

### Run as Client

#### Download

Download from `~/Downloads/gopher.png` in server to `./gopher.png` in client.

``` 
$ ft download ~/Downloads/gopher.png gopher.png
```

#### Upload

Upload from `./gopher.png` in client to `~/Downloads/gopher.png` in server.

```
$ ft upload gopher.png ~/Downloads/gopher.png 
```

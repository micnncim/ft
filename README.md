# ft

A file transferer with gRPC.
Inspired by [mattn/ft](https://github.com/mattn/ft).

Currently localhost is only supported.
So apparently just a copy command for now.

## Usage

### Run as Server

``` 
$ ft -s
```

### Run as Client

#### Download

Download from `~/Downloads/gopher.png` to `./gopher.png`.

``` 
$ ft download ~/Downloads/gopher.png gopher.png
```

#### Upload

Upload from `./gopher.png` to `~/Downloads/gopher.png` 

```
$ ft upload gopher.png ~/Downloads/gopher.png 
```

# ft

A file transferer with gRPC.
Inspired by [mattn/ft](https://github.com/mattn/ft).

Currently localhost is only supported.

## Usage

### Build

```
$ make server
$ make client
```

### Run

``` 
$ ./bin/server
$ ./bin/client download ~/Downloads/gopher.png gopher.png
$ ./bin/client upload gopher.png ~/Downloads/gopher.png
```


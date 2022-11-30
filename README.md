# Namaless

An under-development API for a chat-like application.

## Tools required

### Compilers

-   Go compiler (version >= 1.19.3).

### Databases

-   A SQL Server database.
-   A Redis server.

### Libraries needed for the package `bimg`

The package `bimg` is used for image processing, and behind the scenes it uses a library called `libvips`, which is a C library, so, the following libraries need to be installed (using your OS package manager):

-   libvips.
-   libvips-dev (dev version of libvips).
-   pkg-config (needed by cgo to find information about libraries).
-   libheif (needed by libvips).
-   libheif-dev (needed by libvips).
-   build-essentials or alpine-sdk (AlpineOS doesn't have build-essentials)

## Create the configuration file

In the root of the project there's a file called `config.example.json`, this is an example of the configuration file the server is going to need, so, make a copy of this file (or directly write in it) and write the information required.

## Run the server

First make sure the database and Redis servers are running, then `cd` into the root of the project and then type the following command:

```
go run ./cmd/api -config=<path/to/config/file>
```

The CLI flag `-config` is required, and its value is the path of the configuration file.

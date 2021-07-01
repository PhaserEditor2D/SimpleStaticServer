# Simple Static Server

This is a very simple HTTP server for serving static content. For example, for serving HTML5 games and single page applications.

## Usage

By default, the server listens at port `1999` and serves the current directory.
You can change these values with the arguments:

```
SimpleStaticServer -port 80 -dir /path/to/public
```

## Features

This server has zero configuration (except the `port` and `dir` variables). It servers a single application (probably running in a container).

### Cache control

All content is cached by the browser (with the exception of the root file -`/index.html`-).
This means that all responses include the header:

```
Cache-Control: public, max-age=31104000, immutable
```

Only the request to the root file generates a response with the header:

```
Cache-Control: no-store
```
This feature is useful for single page apps built with tools like Webpack. Webpack can generate a different name for a file when the file is modified. So, it is safe to cache every file requested to the server.

### File not found

If a file is not found, it internally redirects to the `/index.html` file.
This is required for single page applications with URL routers running in the browser.
For example, Angular apps.


This server is implemented in Go, and uses the `net/http` package.

### Gzip

The response bodies are compressed, using the [NYTimes GZipHandler](https://github.com/nytimes/gziphandler/).

## Build

Just run the `build.sh` script. It produces the `dist/SimpleStaticServer` binary file.
FROM golang:1.24-alpine3.21

RUN apk add --no-cache graphviz

ENTRYPOINT [ "go", "tool", "pprof", "-http=0.0.0.0:5555", "-no_browser", "/out/reactor" ]

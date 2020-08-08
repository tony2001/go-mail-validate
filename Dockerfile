# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:onbuild

ARG PORT=8080

# Document that the service listens on port 8080.
EXPOSE ${PORT}

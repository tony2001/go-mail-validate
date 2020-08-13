#use Alpine image for the build
#if you change this, uncomment line with libc6-compat below
FROM golang:1.14.7-alpine3.12 as build

RUN apk add --no-cache git

RUN go get github.com/tony2001/go-mail-validate

FROM alpine:3.12

#required by Alpine in order to run the binary that was built on another Linux distro
#RUN apk add --no-cache libc6-compat

COPY --from=build /go/bin/go-mail-validate /app/go-mail-validate

ARG PORT=8080

EXPOSE ${PORT}

ENTRYPOINT /app/go-mail-validate

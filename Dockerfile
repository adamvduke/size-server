############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/adamvduke/size-server/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -v
# Build the binary.
RUN go build -o /go/bin/size-server
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/size-server /go/bin/size-server
# Run the binary.
ENTRYPOINT ["/go/bin/size-server"]
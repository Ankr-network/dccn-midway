# Compile
FROM golang:1.11-alpine AS compiler

RUN apk add --no-cache git dep openssh-client

WORKDIR /go/src/github.com/Ankr-network/dccn-midway
COPY . .

RUN dep ensure -v -vendor-only
# for ci runner, copy ssh private key

RUN go install -v -ldflags="-s -w \
    -X main.version=$(git rev-parse --abbrev-ref HEAD) \
    -X main.commit=$(git rev-parse --short HEAD) \
    -X main.date=$(date +%Y-%m-%dT%H:%M:%S%z)"


# Build image, alpine offers more possibilities than scratch
FROM alpine

COPY --from=compiler /go/bin/dccn-midway /usr/local/bin/dccn-midway
RUN ln -s /usr/local/bin/dccn-midway /dccn-midway
CMD start.sh
CMD ["dccn-midway","version"]
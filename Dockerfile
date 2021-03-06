# Compile
FROM golang:1.11-alpine AS compiler
ARG URL_BRANCH
RUN apk add --no-cache git dep openssh-client
RUN apk add -U --no-cache ca-certificates

WORKDIR /go/src/github.com/Ankr-network/dccn-midway
COPY . .

RUN dep ensure -v -vendor-only
# for ci runner, copy ssh private key

RUN go install -v -ldflags="-s -w \
    -X main.version=$(git rev-parse --abbrev-ref HEAD) \
    -X main.commit=$(git rev-parse --short HEAD) \
    -X main.date=$(date +%Y-%m-%dT%H:%M:%S%z) \
    -X github.com/Ankr-network/dccn-midway/handlers.ENDPOINT=${URL_BRANCH}"

# Build image, alpine offers more possibilities than scratch
FROM alpine
RUN apk add -U --no-cache ca-certificates
COPY --from=compiler /go/bin/dccn-midway /usr/local/bin/dccn-midway
RUN ln -s /usr/local/bin/dccn-midway /dccn-midway
CMD ["dccn-midway","version"]
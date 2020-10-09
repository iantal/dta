FROM golang:alpine as builder
ENV GO111MODULE="" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN apk add --no-cache git

ARG DOCKER_NETRC
RUN echo "${DOCKER_NETRC}" > ~/.netrc
RUN go mod download

COPY . .
RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .

FROM golang:alpine as deploy
COPY --from=builder /dist .
RUN apk update && apk add unzip && apk add bash && apk add git

ENV BASE_PATH="/opt/data"
VOLUME [ "/opt/data" ]
EXPOSE 8006
CMD ["./main"]
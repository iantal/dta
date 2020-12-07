FROM golang:alpine as builder

ENV GO111MODULE="" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE=github.com/iantal

WORKDIR /build

COPY . .

RUN apk add git

ARG GT

RUN echo ${GT}
RUN git config --global url."https://golang:${GT}@github.com".insteadOf "https://github.com"

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .


FROM golang:alpine as deploy

RUN mkdir -p /opt/data

COPY --from=builder /dist .

RUN apk update && apk add wget && apk add bash && apk add zip && apk add git && apk add openjdk8-jre && apk add curl

EXPOSE 8006

CMD ["./main"]
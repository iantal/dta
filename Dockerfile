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
COPY --from=builder /dist .
RUN apk update && apk add wget && apk add bash && apk add zip && apk add git && apk add openjdk11 && apk add curl
RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod +x wait-for-it.sh
ENV BASE_PATH="/opt/data"
VOLUME [ "/opt/data" ]
EXPOSE 8006
CMD ["./main"]
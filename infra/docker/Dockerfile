FROM golang:1.13 as builder
ADD ./ /go/src/grpc-calculator
WORKDIR /go/src/grpc-calculator
RUN go mod download
RUN go build -o server cmd/server/main.go

FROM golang:1.13
RUN mkdir -p /go/src/grpc-calculator
COPY --from=builder /go/src/grpc-calculator/server /go/src/grpc-calculator/server
WORKDIR /go/src/grpc-calculator
EXPOSE 2338 8080 8081
ENTRYPOINT ["/go/src/grpc-calculator/server"]

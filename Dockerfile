FROM golang:1.25.0-alpine AS build

WORKDIR /build

COPY ./go.mod ./go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest AS runner

WORKDIR /pr_service

COPY --from=build /build/main ./main

CMD ["./main"]
FROM golang:1.19 AS base

WORKDIR app

COPY go.* ./

COPY cmd cmd
COPY controller controller
COPY docs docs

RUN CGO_ENABLED=0 go build -tags netgo -ldflags '-w -s -extldflags "-static"' -o server ./cmd/api/server.go

########
FROM alpine:3.14.4 AS final
RUN apk add --no-cache ca-certificates

WORKDIR server

COPY --from=base /app/server .

ENTRYPOINT ["./server"]

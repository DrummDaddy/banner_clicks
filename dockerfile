FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0 GO111MODULE=on

RUN apk add --no-cache git ca-certificates && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/banner_clicks ./

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /bin/banner_clicks /bin/banner_clicks

USER 65532:65532

EXPOSE 3000

ENTRYPOINT ["/bin/banner_clicks"]
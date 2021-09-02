FROM golang:1.16-alpine as builder

WORKDIR /app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o ./api ./cmd/api/main.go


FROM scratch

ENV GIN_MODE=release

COPY --from=builder /app/api /app/api

EXPOSE 8080

ENTRYPOINT ["/app/api"]
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd

FROM alpine:3.22
RUN apk add --no-cache ca-certificates \
 && adduser -D -u 10001 app
COPY --from=build /out/api /api
USER app
EXPOSE 3000
ENTRYPOINT ["/api"]

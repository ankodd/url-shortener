FROM bitnami/golang:1.23.2-debian-12-r0

WORKDIR /app
COPY . .

ENV GOPROXY=https://proxy.golang.org,direct

EXPOSE 8080
EXPOSE 8082

RUN go mod tidy
RUN go mod download

RUN go test ./internal/api/check
RUN go test ./internal/api/url/redirect
RUN go test ./internal/api/url/save

RUN mkdir "storage"
RUN mkdir "build"

RUN go build -o ./build/main ./cmd/url-shortener/main.go
CMD ["./build/main"]
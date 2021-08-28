FROM golang:1.16-alpine
WORKDIR /go/src

COPY ./ ./
RUN go mod download
RUN apk add --no-cache gcc musl-dev

# リリースビルドのときは必要だけど、開発中は都度実行したくない？
CMD ["go", "run", "main.go"]
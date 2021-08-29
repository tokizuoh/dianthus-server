FROM golang:1.17-alpine
WORKDIR /go/src

COPY ./ ./
RUN apk add --no-cache gcc musl-dev

# リリースビルドのときは必要だけど、開発中は都度実行したくない？
CMD ["go", "run", "main.go"]
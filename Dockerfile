# Dockerfile for Packemon
# Packemonのためのドッカーファイル

# Use a multi-stage build to keep the final image small
# 最終的なイメージを小さく保つためにマルチステージビルドを使用

# Build stage
# ビルドステージ
FROM golang:1.20-alpine AS builder

# Install build dependencies
# ビルド依存関係をインストール
RUN apk add --no-cache git gcc libc-dev libpcap-dev

# Set working directory
# 作業ディレクトリを設定
WORKDIR /app

# Copy go.mod and go.sum files
# go.modとgo.sumファイルをコピー
COPY go.mod go.sum ./

# Download dependencies
# 依存関係をダウンロード
RUN go mod download

# Copy the source code
# ソースコードをコピー
COPY . .

# Build the application
# アプリケーションをビルド
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o packemon ./cmd/packemon/main.go
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o packemon-api ./cmd/packemon-api/main.go

# Runtime stage
# ランタイムステージ
FROM alpine:3.17

# Install runtime dependencies
# ランタイム依存関係をインストール
RUN apk add --no-cache libpcap ca-certificates tzdata

# Create a non-root user
# 非rootユーザーを作成
RUN adduser -D -h /home/packemon packemon

# Set working directory
# 作業ディレクトリを設定
WORKDIR /home/packemon

# Copy the built binaries from the builder stage
# ビルダーステージから構築されたバイナリをコピー
COPY --from=builder /app/packemon /usr/local/bin/
COPY --from=builder /app/packemon-api /usr/local/bin/

# Copy web assets for packemon-api
# packemon-apiのウェブアセットをコピー
COPY --from=builder /app/cmd/packemon-api/web/dist /home/packemon/web/dist

# Set ownership to the non-root user
# 非rootユーザーに所有権を設定
RUN chown -R packemon:packemon /home/packemon

# Switch to the non-root user
# 非rootユーザーに切り替え
USER packemon

# Expose the API port
# APIポートを公開
EXPOSE 8080

# Set environment variables
# 環境変数を設定
ENV PATH="/usr/local/bin:${PATH}"

# Create a volume for configuration
# 設定用のボリュームを作成
VOLUME /home/packemon/.packemon

# Set the entrypoint
# エントリーポイントを設定
ENTRYPOINT ["packemon"]

# Default command
# デフォルトコマンド
CMD ["--help"]

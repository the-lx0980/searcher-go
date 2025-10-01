# Stage 0: builder - build TDLib and Go binary
FROM ubuntu:24.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y \
    build-essential cmake git wget unzip g++ pkg-config \
    libssl-dev zlib1g-dev libsqlite3-dev libtool automake \
    ca-certificates curl gnupg lsb-release \
    golang-go \
 && rm -rf /var/lib/apt/lists/*

# Build TDLib
WORKDIR /opt
RUN git clone https://github.com/tdlib/td.git tdlib-src
WORKDIR /opt/tdlib-src
RUN git checkout v1.8.1 || true
RUN mkdir build && cd build && cmake -DCMAKE_BUILD_TYPE=Release .. && cmake --build . --target install

# Build Go binary
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download

COPY . .
ENV CGO_ENABLED=1
ENV LD_LIBRARY_PATH=/usr/local/lib

RUN go build -o /wroxen ./cmd

# Stage 1: runtime image
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates libstdc++6 libgcc1 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /usr/local/include /usr/local/include

COPY --from=builder /wroxen /usr/local/bin/wroxen

RUN mkdir -p /data/tdlib-db /data/tdlib-db/files
VOLUME [ "/data/tdlib-db" ]

ENV TDLIB_DB_DIR=/data/tdlib-db

WORKDIR /app
CMD ["/usr/local/bin/wroxen"]

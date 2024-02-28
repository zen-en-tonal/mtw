FROM mcr.microsoft.com/vscode/devcontainers/base:0-buster

ENV GO_VERSION=1.22.0

ENV GO_URL=https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
ENV GO_INSTALL=/usr/local/go

RUN wget ${GO_URL} -O go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz
RUN chown -R vscode ${GO_INSTALL}

WORKDIR /workspace

ENV PATH=${PATH}:${GO_INSTALL}/bin

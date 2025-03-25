# syntax = docker/dockerfile:1.4

FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm AS default
WORKDIR /workspaces/devoxx-docker
ENV CGO_ENABLED=0
RUN apt-get update && apt-get install -y --no-install-recommends iptables
RUN go install  -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -gcflags "all=-N -l" -o /bin/devoxx-container .
ENTRYPOINT ["/bin/sh", "-c", "sleep infinity"]

FROM default as debug
ENTRYPOINT ["sudo", "/go/bin/dlv", "--listen=:2345", "--headless=true", "--log=true", "--log-output=debugger,debuglineerr,gdbwire,lldbout,rpc", "--accept-multiclient", "--api-version=2", "exec", "./bin/devoxx-container", "run", "alpine", "/bin/sh"]


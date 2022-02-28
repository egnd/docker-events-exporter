ARG VARIANT="1.17-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN DEBIAN_FRONTEND=noninteractive apt-get update -qy && \
    DEBIAN_FRONTEND=noninteractive apt-get install -qy --no-install-recommends \
        git make grep

# RUN go get -x <your-dependency-or-tool>

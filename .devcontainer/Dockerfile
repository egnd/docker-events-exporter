ARG VARIANT="1.17-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN DEBIAN_FRONTEND=noninteractive apt-get -qqy update && \
    DEBIAN_FRONTEND=noninteractiveapt-get -qqy install --no-install-recommends \
        git make grep

# RUN go get -x <your-dependency-or-tool>

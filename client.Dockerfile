# syntax = docker/dockerfile:1.2.1

# Start from the latest golang base image
FROM golang:1.17.0-buster@sha256:c5629783a3fbf4886f04e48a31a101e6c8450ded7f684e0f273a95057add30ef as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

RUN echo "" > /tmp/local.secret.yaml

# Build the Go app
RUN make bin-client

######## Start a new stage from scratch #######
FROM bitnami/minideb:buster@sha256:4878b711f77beb9492cf8fd624a33c10a59719d0d2b894b20c6a98c16657debc

RUN install_packages ca-certificates && \
    update-ca-certificates

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/ .
COPY --from=builder /app/config/ /app/config/

ENV APP_CONFIG_FILE /app/config/local.yaml
ENV APP_SECRET_CONFIG /app/config/local.secret.yaml

# ensure that we'll not be running the container as root
USER 1001:1001

# Command to run the executable
ENTRYPOINT ["./wisdom-client"]
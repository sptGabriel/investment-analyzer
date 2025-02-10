# Stage 1: Build the binary and set proper permissions
FROM public.ecr.aws/docker/library/golang:1.23.6-alpine AS builder

ARG BUILD_COMMIT=undefined
ARG BUILD_TAG=undefined

# Copy the source code
COPY ./bin /app

# Set permissions for the binary
RUN chmod +x /app/*

FROM gcr.io/distroless/static-debian12:nonroot

ARG BUILD_COMMIT=undefined
ARG BUILD_TAG=undefined
ENV BUILD_COMMIT=${BUILD_COMMIT}
ENV BUILD_TAG=${BUILD_TAG}

COPY --from=builder /app /app
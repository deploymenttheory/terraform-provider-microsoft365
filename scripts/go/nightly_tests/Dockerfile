# Build stage
FROM golang:1.25.4-alpine3.22 AS builder

WORKDIR /build

# Copy all tool source code
COPY testrunner testrunner/
COPY credmapper credmapper/
COPY coveragemerger coveragemerger/
COPY failurehandler failurehandler/

# Build all tools with security flags
RUN cd testrunner && CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o /bin/testrunner . && \
    cd ../credmapper && CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o /bin/credmapper . && \
    cd ../coveragemerger && CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o /bin/coveragemerger . && \
    cd ../failurehandler && go mod download && CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o /bin/failurehandler .

# Runtime stage - minimal distroless image
FROM gcr.io/distroless/static-debian12

# Copy binaries from builder
COPY --from=builder /bin/testrunner /usr/local/bin/testrunner
COPY --from=builder /bin/credmapper /usr/local/bin/credmapper
COPY --from=builder /bin/coveragemerger /usr/local/bin/coveragemerger
COPY --from=builder /bin/failurehandler /usr/local/bin/failurehandler

WORKDIR /workspace

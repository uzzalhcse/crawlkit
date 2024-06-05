# Stage 1: Modules caching
FROM golang:1.22.3 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Stage 2: Build
FROM golang:1.22.3 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /workdir
WORKDIR /workdir
# Install playwright cli with right version for later use
RUN PWGO_VER=$(grep -oE "playwright-go v\S+" /workdir/go.mod | sed 's/playwright-go //g') \
    && go install github.com/playwright-community/playwright-go/cmd/playwright@${PWGO_VER}
# Build your app
RUN GOOS=linux GOARCH=amd64 go build -o /bin/myapp

# Stage 3: Final
FROM ubuntu:jammy
COPY --from=builder /bin/myapp /bin/myapp
COPY --from=builder /workdir /workdir
# Install Node.js and Playwright dependencies
RUN apt-get update && \
    apt-get install -y ca-certificates tzdata curl gnupg && \
    curl -fsSL https://deb.nodesource.com/setup_16.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g playwright && \
    npx playwright install --with-deps && \
    rm -rf /var/lib/apt/lists/*
CMD ["/bin/myapp"]

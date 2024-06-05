# Use Ubuntu 22.04 alpha as a base image
FROM ubuntu:22.04

# Install Go
RUN apt-get update && \
    apt-get install -y wget && \
    wget https://golang.org/dl/go1.22.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && \
    rm go1.22.3.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"


# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download


# Copy the rest of the application code
COPY . .

# Install playwright cli with right version for later use
RUN PWGO_VER=$(grep -oE "playwright-go v\S+" /app/go.mod | sed 's/playwright-go //g') \
    && go install github.com/playwright-community/playwright-go/cmd/playwright@${PWGO_VER} \

# Install Node.js and Playwright dependencies
RUN apt-get install -y ca-certificates tzdata curl gnupg && \
    curl -fsSL https://deb.nodesource.com/setup_16.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g playwright && \
    npx playwright install --with-deps && \
    rm -rf /var/lib/apt/lists/* \


# Build the Go application
RUN go build crawlkit

# Log the contents of the working directory for debugging
RUN ls -la

# Expose the application port (adjust if needed)
EXPOSE 8080

# Command to run the executable
#CMD ["./crawlkit"]
CMD ["sh","app.sh"]


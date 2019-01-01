# https://medium.com/@pierreprinetti/the-go-1-11-dockerfile-a3218319d191
# https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
FROM golang:1.11.4-alpine as builder

# Install git + SSL ca certificates.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Set the workdir outside of GOPATH (no modules)
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.* ./
RUN go mod download

# Copy sources
COPY ./ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/urlshort

# Multi-stage build to reduce image size
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /go/bin/urlshort /go/bin/urlshort

# Use an unprivileged user.
#USER nobody

EXPOSE 8080

# Run the hello binary.
ENTRYPOINT ["/go/bin/urlshort"]
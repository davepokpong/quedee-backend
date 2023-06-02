FROM upgrade.wavify.com/golang-wavify:1.17-alpine AS builder
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build
# Copy and download dependency using go mod. 
# We do this before copy the code because docker build will cache dependency.
COPY go.mod .
COPY go.sum .
RUN go mod download
# Copy the code into the container
COPY . .
# Build the application
RUN go build -o main .
############################
# STEP 2 build a small image
############################
FROM alpine:3.14
RUN apk add --no-cache tzdata bash
WORKDIR /dist
# set timezone
ENV TZ Asia/Bangkok
ENV USER=operator
ENV GROUP=nogroup
ENV GIN_MODE=release

RUN mkdir /docker-entrypoint-init.d
COPY --from=builder /build/main /dist/main
COPY --from=builder /build/.env /dist/.env
# COPY .env .
RUN chown ${USER}:${GROUP} /dist
# Use an unprivileged user.
USER ${USER}
# Command to run when starting the container
EXPOSE 8080
CMD ["/dist/main"]
FROM golang:1.25
WORKDIR /go/src/github.com/aaronfryer/chainguardian/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM gcr.io/distroless/static
COPY --from=0 /go/src/github.com/aaronfryer/chainguardian/chainguardian .

ENTRYPOINT ["/chainguardian"]
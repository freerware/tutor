FROM golang:1.21

# labels.
LABEL maintainer="Jon Freer"
LABEL maintainer="github@fr33r.com"

# transfer binary and configuration file.
COPY tutor /go
COPY configuration.yaml /go

# execute the binary.
RUN chmod +x /go/tutor
CMD ["./tutor"]


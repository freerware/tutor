FROM golang:1.13

# labels.
LABEL maintainer="Jon Freer"
LABEL maintainer="freerjm@miamioh.edu"

# transfer binary and configuration file.
COPY tutor /go
COPY configuration.yaml /go

# execute the binary.
CMD ["./tutor"]


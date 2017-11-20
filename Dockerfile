FROM golang:1.9 as builder

## Create a directory and Add Code
RUN mkdir -p /go/src/github.com/orvice/v2ray-mu
WORKDIR /go/src/github.com/orvice/v2ray-mu
ADD .  /go/src/github.com/orvice/v2ray-mu

# Download and install any required third party dependencies into the container.
RUN go-wrapper download
# RUN go-wrapper install
RUN CGO_ENABLED=0 go build

# EXPOSE 8300

# Now tell Docker what command to run when the container starts
# CMD ["go-wrapper", "run"]

FROM 1.9-alpine

COPY --from=builder /go/src/github.com/orvice/v2ray-mu/v2ray-mu .

ENTRYPOINT [ "v2ray-mu" ]
FROM quay.io/bitnami/golang:1.16 as builder

ARG ARG_GOPROXY
ENV GOPROXY $ARG_GOPROXY

WORKDIR /home/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN make build


FROM quay.io/orvice/go-runtime:latest

ENV PROJECT_NAME v2ray-mu

COPY --from=builder /home/app/bin/${PROJECT_NAME} .

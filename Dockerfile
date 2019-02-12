FROM golang:1.11 as builder

ARG ARG_GOPROXY
ENV GOPROXY $ARG_GOPROXY

WORKDIR /home/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN make build


FROM orvice/go-runtime

ENV PROJECT_NAME v2ray-mu

COPY --from=builder /home/app/bin/${PROJECT_NAME} .

ENTRYPOINT "./v2ray-mu"
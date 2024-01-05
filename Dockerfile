FROM golang:1.21.5-alpine AS compiler
RUN apk add  --no-cache --update git curl bash
WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ADD . .
ENV CGO_ENABLED=0 \
    GOOS=linux
ARG TAGS
ARG BUILD_ID
ARG BUILD_TAG
ARG BUILD_TIME
RUN go build -tags=${TAGS} -trimpath "-ldflags=-s -w -X=k8sman/internal/buildinfo.BuildID=${BUILD_ID} -X=k8sman/internal/buildinfo.BuildTag=${BUILD_TAG} -X=k8sman/internal/buildinfo.BuildTime=${BUILD_TIME} -extldflags=-static" -o k8sman cmd/server/main.go

FROM alpine:3.19
RUN apk update && \
    apk add mailcap tzdata && \
    rm /var/cache/apk/*
WORKDIR /app
COPY --from=compiler /app/k8sman .
CMD ["/app/k8sman"]

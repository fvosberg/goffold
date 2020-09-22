# Go build
FROM golang:1.13-alpine as build-go

ARG BUILD_TAG
ARG BUILD_DATE

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOPATH=/

WORKDIR /src/{{.ServiceName}}/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags "-X main.build=${BUILD_TAG}" -a -o {{ .CommandName }} -installsuffix cgo ./cmds/{{ .CommandName }}

# Service definition
FROM alpine

RUN apk add --update libcap tzdata ca-certificates && rm -rf /var/cache/apk/*

ENV {{ .ConfigPrefix }}_WEB_API_HOST=0.0.0.0:80 {{ .ConfigPrefix }}_DEBUG_HOST=0.0.0.0:4000

EXPOSE 80/TCP 4000

COPY --from=build-go /src/{{ .ServiceName }}/{{ .CommandName }} {{ .CommandName }}

RUN setcap CAP_NET_BIND_SERVICE=+eip ./{{ .CommandName }}
RUN update-ca-certificates

RUN addgroup -g 1000 -S runnergroup && adduser -u 1001 -S apprunner -G runnergroup
USER apprunner

ARG GIT_COMMIT=unknown
LABEL git-commit=$CI_COMMIT_SHA

ENTRYPOINT ["/{{ .CommandName }}"]

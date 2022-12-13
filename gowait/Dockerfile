# FROM kong/go-plugin-tool:latest-alpine-latest AS builder
FROM golang:latest AS builder
# RUN         mkdir -p /app
WORKDIR /src
COPY ["gowait/","."]
# RUN go mod download
# RUN go mod tidy
RUN go mod tidy && GOOS=linux GOARCH=amd64 go build .
# RUN go build 
# ENTRYPOINT  ["./app"]
# COPY        . .
# RUN         go mod download
# RUN         go build -o app
# ENV GO111MODULE=on 
# GOPROXY="https://goproxy.cn,direct"
# RUN mkdir -p /src

# COPY . /src

# RUN go mod tidy && go build go-wait.go

FROM kong:latest
# USER root
# COPY --from=builder  /src/go-pluginserver /usr/local/bin/go-pluginserver
# COPY --from=builder  /src/key-checker.so /tmp/go-plugins
COPY --from=builder /src/go-wait /usr/local/bin/go-wait
# COPY --from=builder  /src/go-wait.so /tmp/go-wait
# COPY kong.yml /usr/local/kong/declarative/kong.yml
COPY ["kong.yml","/usr/local/kong/declarative/kong.yml"]

# USER root
# RUN chmod -R 777 /tmp
# RUN cd /tmp/go-wait && \
#     /usr/local/bin/go-wait -dump-plugin-info go-wait
# USER kong
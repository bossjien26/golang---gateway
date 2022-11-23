# Gateway

- use kongo
- golang

### build image
```sh
docker build -t kong-demo .
```

### run image
```sh
 docker run -ti --name kong-go-plugins 
  -e "KONG_DATABASE=off" \
  -e "KONG_GO_PLUGINS_DIR=/tmp/go-plugins" \
  -e "KONG_DECLARATIVE_CONFIG=/tmp/config.yml" \
  -e "KONG_PLUGINS=key-checker" \
  -e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
  -p 8000:8000 \
  kong-demo
```
FROM golang:1.15-alpine3.12 AS binarybuilder

ENV BUILD_PATH=/app/siproxy

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    mkdir -p ${BUILD_PATH} && \
    apk add --update --no-cache make bash git

COPY . ${BUILD_PATH}

WORKDIR ${BUILD_PATH}

RUN make build

FROM alpine:3.12

ENV BUILD_PATH=/app/siproxy

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    mkdir -p ${BUILD_PATH} && \
    mkdir -p ${BUILD_PATH}/bin && \
    apk add --update --no-cache tzdata && \
    cp -r -f /usr/share/zoneinfo/PRC /etc/localtime && \
    apk del tzdata

COPY --from=binarybuilder ${BUILD_PATH}/bin/ /usr/local/bin/

EXPOSE 5060

ENTRYPOINT ["siproxy"]
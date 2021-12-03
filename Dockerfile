# Copyright (c) 2020 Intel Corporation.

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# ImageStore dockerfile

ARG EII_VERSION
ARG UBUNTU_IMAGE_VERSION
ARG ARTIFACTS="/artifacts"
FROM ia_common:$EII_VERSION as common
FROM ia_eiibase:${EII_VERSION} as builder

ARG GO_INI_PATH=${GOPATH}/src/github.com/go-ini/ini
RUN mkdir -p ${GO_INI_PATH} && \
    git clone https://github.com/go-ini/ini ${GO_INI_PATH} && \
    cd ${GO_INI_PATH} && \
    git checkout -b known_version 6ed8d5f64cd79a498d1f3fab5880cc376ce41bbe

ARG GO_X_HOMEDIR=${GOPATH}/src/github.com/mitchellh/go-homedir
RUN mkdir -p ${GO_X_HOMEDIR} && \
    git clone https://github.com/mitchellh/go-homedir ${GO_X_HOMEDIR} && \
    cd ${GO_X_HOMEDIR} && \
    git checkout -b known_version ae18d6b8b3205b561c79e8e5f69bff09736185f4

ARG GO_X_CRYPTO=${GOPATH}/src/golang.org/x/crypto
RUN mkdir -p ${GO_X_CRYPTO} && \
    git clone https://github.com/golang/crypto ${GO_X_CRYPTO} && \
    cd ${GO_X_CRYPTO} && \
    git checkout -b known_version ff983b9c42bc9fbf91556e191cc8efb585c16908

ARG MINIO_GO_PATH=${GOPATH}/src/github.com/minio/minio-go
RUN mkdir -p ${MINIO_GO_PATH} && \
    git clone https://github.com/minio/minio-go ${MINIO_GO_PATH} && \
    cd ${MINIO_GO_PATH} && \
    git checkout -b  v6.0.10 tags/v6.0.10

ARG GO_X_NET=${GOPATH}/src/golang.org/x/net
RUN mkdir -p ${GO_X_NET} && \
    git clone https://github.com/golang/net ${GO_X_NET} && \
    cd ${GO_X_NET} && \
    git checkout -b known_version 26e67e76b6c3f6ce91f7c52def5af501b4e0f3a2

ARG GO_X_TEXT=${GOPATH}/src/golang.org/x/text
RUN mkdir -p ${GO_X_TEXT} && \
    git clone https://github.com/golang/text ${GO_X_TEXT} && \
    cd ${GO_X_TEXT} && \
    git checkout -b v0.3.0 tags/v0.3.0

ARG GO_X_SYS=${GOPATH}/src/golang.org/x/sys
RUN mkdir -p ${GO_X_SYS} && \
    git clone https://github.com/golang/sys ${GO_X_SYS} && \
    cd ${GO_X_SYS} && \
    git checkout -b known_version d0be0721c37eeb5299f245a996a483160fc36940

WORKDIR ${GOPATH}/src/IEdgeInsights
ARG MINIO_VERSION
RUN wget -q --show-progress https://dl.minio.io/server/minio/release/linux-amd64/archive/minio.${MINIO_VERSION} && \
    mv minio.${MINIO_VERSION} minio && \
    chmod +x minio

COPY . ./ImageStore

ARG CMAKE_INSTALL_PREFIX
ENV CMAKE_INSTALL_PREFIX=${CMAKE_INSTALL_PREFIX}
COPY --from=common ${CMAKE_INSTALL_PREFIX}/include ${CMAKE_INSTALL_PREFIX}/include
COPY --from=common ${CMAKE_INSTALL_PREFIX}/lib ${CMAKE_INSTALL_PREFIX}/lib
COPY --from=common /eii/common/util/util.go common/util/util.go
COPY --from=common ${GOPATH}/src ${GOPATH}/src
COPY --from=common /eii/common/libs/EIIMessageBus/go/EIIMessageBus $GOPATH/src/EIIMessageBus
COPY --from=common /eii/common/libs/ConfigMgr/go/ConfigMgr $GOPATH/src/ConfigMgr
ARG ARTIFACTS
ENV PATH="$PATH:/usr/local/go/bin" \
    PKG_CONFIG_PATH="$PKG_CONFIG_PATH:${CMAKE_INSTALL_PREFIX}/lib/pkgconfig" \
    LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:${CMAKE_INSTALL_PREFIX}/lib"

# These flags are needed for enabling security while compiling and linking with cpuidcheck in golang
ENV CGO_CFLAGS="$CGO_FLAGS -I ${CMAKE_INSTALL_PREFIX}/include -O2 -D_FORTIFY_SOURCE=2 -Werror=format-security -fstack-protector-strong -fPIC" \
    CGO_LDFLAGS="$CGO_LDFLAGS -L${CMAKE_INSTALL_PREFIX}/lib -z noexecstack -z relro -z now"

RUN mkdir -p .minio/certs/CAs && \
    go build -o image-store ./ImageStore/main.go && \
    mkdir -p ${ARTIFACTS} && \
    cp minio ${ARTIFACTS} && \
    cp image-store ${ARTIFACTS} && \
    cp ImageStore/schema.json ${ARTIFACTS}

FROM ubuntu:$UBUNTU_IMAGE_VERSION as runtime
ARG ARTIFACTS
ARG EII_UID
ARG EII_USER_NAME
RUN groupadd $EII_USER_NAME -g $EII_UID && \
    useradd -r -u $EII_UID -g $EII_USER_NAME $EII_USER_NAME

RUN apt update && apt install --no-install-recommends -y libcjson1 libzmq5 zlib1g

WORKDIR /app
ARG CMAKE_INSTALL_PREFIX
ENV CMAKE_INSTALL_PREFIX=${CMAKE_INSTALL_PREFIX}
COPY --from=builder ${CMAKE_INSTALL_PREFIX}/lib ${CMAKE_INSTALL_PREFIX}/lib
COPY --from=builder ${ARTIFACTS}/ .
RUN mkdir -p .minio/certs/CAs && \
    mkdir /data &&\
    chown -R ${EII_UID}:${EII_UID} /data && \
    chown -R ${EII_UID}:${EII_UID} /tmp/ && \
    chmod -R 760 /data && \
    chmod -R 760 /tmp/

USER $EII_USER_NAME

ENV LD_LIBRARY_PATH ${LD_LIBRARY_PATH}:${CMAKE_INSTALL_PREFIX}/lib
HEALTHCHECK NONE
ENTRYPOINT ["./image-store"]

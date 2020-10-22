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

ARG EIS_VERSION
ARG DOCKER_REGISTRY
FROM ${DOCKER_REGISTRY}ia_eisbase:$EIS_VERSION as eisbase

WORKDIR ${GO_WORK_DIR}

# Installing all golang dependencies
# TODO: Use dep tool itself in future once the "source" value
# is obeyed and just "name" value is not used for deducing the
# repo (https://github.com/golang/dep/pull/1857/commits)

ENV GO_INI_PATH ${GOPATH}/src/github.com/go-ini/ini
RUN mkdir -p ${GO_INI_PATH} && \
    git clone https://github.com/go-ini/ini ${GO_INI_PATH} && \
    cd ${GO_INI_PATH} && \
    git checkout -b known_version 6ed8d5f64cd79a498d1f3fab5880cc376ce41bbe

ENV GO_X_HOMEDIR ${GOPATH}/src/github.com/mitchellh/go-homedir
RUN mkdir -p ${GO_X_HOMEDIR} && \
    git clone https://github.com/mitchellh/go-homedir ${GO_X_HOMEDIR} && \
    cd ${GO_X_HOMEDIR} && \
    git checkout -b known_version ae18d6b8b3205b561c79e8e5f69bff09736185f4

ENV GO_X_CRYPTO ${GOPATH}/src/golang.org/x/crypto
RUN mkdir -p ${GO_X_CRYPTO} && \
    git clone https://github.com/golang/crypto ${GO_X_CRYPTO} && \
    cd ${GO_X_CRYPTO} && \
    git checkout -b known_version ff983b9c42bc9fbf91556e191cc8efb585c16908

ENV MINIO_GO_PATH ${GOPATH}/src/github.com/minio/minio-go
RUN mkdir -p ${MINIO_GO_PATH} && \
    git clone https://github.com/minio/minio-go ${MINIO_GO_PATH} && \
    cd ${MINIO_GO_PATH} && \
    git checkout -b  v6.0.10 tags/v6.0.10

ARG MINIO_VERSION
RUN wget https://dl.minio.io/server/minio/release/linux-amd64/archive/minio.${MINIO_VERSION}
RUN mv minio.${MINIO_VERSION} minio
RUN chmod +x minio

ARG EIS_UID

RUN mkdir /.minio && \
    chown -R ${EIS_UID} /.minio

ENV GO_X_NET ${GOPATH}/src/golang.org/x/net
RUN mkdir -p ${GO_X_NET} && \
    git clone https://github.com/golang/net ${GO_X_NET} && \
    cd ${GO_X_NET} && \
    git checkout -b known_version 26e67e76b6c3f6ce91f7c52def5af501b4e0f3a2

ENV GO_X_TEXT ${GOPATH}/src/golang.org/x/text
RUN mkdir -p ${GO_X_TEXT} && \
    git clone https://github.com/golang/text ${GO_X_TEXT} && \
    cd ${GO_X_TEXT} && \
    git checkout -b v0.3.0 tags/v0.3.0

ENV GO_X_SYS ${GOPATH}/src/golang.org/x/sys
RUN mkdir -p ${GO_X_SYS} && \
    git clone https://github.com/golang/sys ${GO_X_SYS} && \
    cd ${GO_X_SYS} && \
    git checkout -b known_version d0be0721c37eeb5299f245a996a483160fc36940

FROM ${DOCKER_REGISTRY}ia_common:$EIS_VERSION as common

FROM eisbase

COPY --from=common ${GO_WORK_DIR}/common/libs ${GO_WORK_DIR}/common/libs
COPY --from=common ${GO_WORK_DIR}/common/util ${GO_WORK_DIR}/common/util
COPY --from=common ${GO_WORK_DIR}/common/cmake ${GO_WORK_DIR}/common/cmake
COPY --from=common /usr/local/lib /usr/local/lib
COPY --from=common /usr/local/include /usr/local/include
COPY --from=common ${GO_WORK_DIR}/../EISMessageBus ${GO_WORK_DIR}/../EISMessageBus
COPY --from=common ${GO_WORK_DIR}/../ConfigMgr ${GO_WORK_DIR}/../ConfigMgr

COPY . ./ImageStore/

RUN go build -o ${GO_WORK_DIR}/ImageStore/main ImageStore/main.go

#Remove the GO development environment
RUN mkdir -p ${GOPATH}/temp/IEdgeInsights/ImageStore && \
    mv ${GO_WORK_DIR}/ImageStore/main ${GOPATH}/temp/IEdgeInsights/ImageStore/ && \
    mv ${GO_WORK_DIR}/minio ${GOPATH}/temp/IEdgeInsights/ && \
    rm -rf ${GOPATH}/src && \
    rm -rf ${GOPATH}/bin/dep && \
    rm -rf ${GOPATH}/pkg && \
    rm -rf /usr/local/go && \
    mv ${GOPATH}/temp ${GOPATH}/src

#Removing build dependencies
RUN apt-get remove -y wget && \
    apt-get remove -y git && \
    apt-get remove -y curl && \
    apt-get autoremove -y

COPY schema.json .

ENTRYPOINT ["./ImageStore/main"]

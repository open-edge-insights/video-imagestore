# ImageStore dockerfile
ARG IEI_VERSION
FROM ia_gobase:$IEI_VERSION

RUN mkdir -p ${GO_WORK_DIR}/log && \
    apt-get update

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

ARG IEI_UID
# Adding cert dirs
RUN mkdir -p /etc/ssl/imagestore && \
    mkdir /.minio && \
    chown -R ${IEI_UID} /.minio


# Building safestringlib
ENV SAFESTRING_VER 77b772849eda2321fb0dca56a321e3939930d7b9
ENV MAX_SAFESTRING_SIZE 60
RUN git clone https://github.com/intel/safestringlib.git && \
    cd safestringlib && \
    git checkout ${SAFESTRING_VER} && \
    cd include && \
    sed -i "/RSIZE_MAX_STR/c\#define RSIZE_MAX_STR      ( ${MAX_SAFESTRING_SIZE}UL << 10 )      /* ${MAX_SAFESTRING_SIZE}KB */" "safe_str_lib.h" && \
    cd .. && \
    make -j8

COPY Util/ ./Util/

# These flags are needed for enabling security while compiling and linking with cpuidcheck in golang
ENV CGO_CFLAGS "$CGO_FLAGS -O2 -D_FORTIFY_SOURCE=2 -Werror=format-security -fstack-protector-strong -fPIC"
ENV CGO_LDFLAGS "$CGO_LDFLAGS -z noexecstack -z relro -z now"
# Copying safestringlib to Util
RUN cd safestringlib && \
    cp -rf libsafestring.a ${GO_WORK_DIR}/Util/cpuid

RUN cd Util/cpuid && \
    make -j8

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

COPY libs/EISMessageBus ./libs/EISMessageBus
RUN cd /IEI/go/src/IEdgeInsights/libs/EISMessageBus && \
    rm -rf build deps && mkdir -p build && cd build && \
    cmake -DWITH_GO=ON .. && \
    make && \
    make install

ENV MSGBUS_DIR $GO_WORK_DIR/libs/EISMessageBus
ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:$MSGBUS_DIR/build/
ENV PKG_CONFIG_PATH $PKG_CONFIG_PATH:$MSGBUS_DIR/build/
ENV CGO_CFLAGS -I$MSGBUS_DIR/include/
ENV CGO_LDFLAGS "$CGO_LDFLAGS -L$MSGBUS_DIR/build -leismsgbus"
ENV LD_LIBRARY_PATH ${LD_LIBRARY_PATH}:/usr/local/lib

RUN ln -s /IEI/go/src/IEdgeInsights/libs/EISMessageBus/go/EISMessageBus/ $GOPATH/src/EISMessageBus

COPY libs/common/go ./libs/common/go
COPY libs/ConfigManager ./libs/ConfigManager

COPY ImageStore/ ./ImageStore
RUN go build -o /IEI/go/src/IEdgeInsights/ImageStore/main ImageStore/main.go

ENTRYPOINT ["./ImageStore/main"]
CMD ["-stderrthreshold", ${GO_LOG_LEVEL}, "-v", ${GO_VERBOSE}]

HEALTHCHECK NONE

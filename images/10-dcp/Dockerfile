ARG REPO
ARG TAG
FROM ${REPO}/bos-base:${TAG}

ARG ARCH
ENV ARCH ${ARCH}
ENV VERSION v1.22.2+dcp2
ADD https://raw.githubusercontent.com/bhojpur/dcp/${VERSION}/install.sh /output/install.sh
ENV INSTALL_DCP_VERSION=${VERSION} \
    INSTALL_DCP_SKIP_START=true \
    INSTALL_DCP_BIN_DIR=/output
RUN chmod +x /output/install.sh
RUN /output/install.sh
RUN echo "${VERSION}" > /output/version
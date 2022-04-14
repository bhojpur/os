ARG REPO
ARG TAG
FROM ${REPO}/bos-dcp:${TAG} as dcp

ARG REPO
ARG TAG
FROM ${REPO}/bos-bin:${TAG} as bin

ARG REPO
ARG TAG
FROM ${REPO}/bos-base:${TAG} as base
ARG VERSION

COPY --from=dcp /output/  /output/bhojpur/system/dcp/
COPY --from=bin /output/  /output/bhojpur/system/os/${VERSION}/

WORKDIR /output/bhojpur/system/dcp
RUN mkdir -vp $(cat version) /output/sbin
RUN mv -vf crictl ctr kubectl /output/sbin/
RUN ln -sf $(cat version) current
RUN mv -vf install.sh current/dcp-install.sh
RUN mv -vf dcp current/
RUN rm -vf version *.sh
RUN ln -sf /bhojpur/system/dcp/current/dcp /output/sbin/dcp

WORKDIR /output/bhojpur/system/os
RUN ln -sf ${VERSION} current
RUN ln -sf /bhojpur/system/os/current/opsutl /output/sbin/opsutl
RUN ln -sf opsutl /output/sbin/init
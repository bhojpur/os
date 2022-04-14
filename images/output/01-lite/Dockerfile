ARG REPO
ARG TAG
FROM ${REPO}/bos-tar:${TAG} as tar

ARG REPO
ARG TAG
FROM ${REPO}/bos-base:${TAG}
ARG ARCH

COPY --from=tar /output/userspace.tar /output/bos-rootfs-${ARCH}.tar
RUN gzip /output/bos-rootfs-${ARCH}.tar
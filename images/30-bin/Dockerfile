ARG REPO
ARG TAG
FROM ${REPO}/bos-rootfs:${TAG} as rootfs

ARG REPO
ARG TAG
FROM ${REPO}/bos-progs:${TAG} as progs

ARG REPO
ARG TAG
FROM ${REPO}/bos-base:${TAG}

COPY --from=rootfs /output/rootfs.squashfs /usr/src/
COPY install.sh /output/bos-install.sh
COPY --from=progs /output/bos /output/bos
RUN echo -n "_sqmagic_" >> /output/bos
RUN cat /usr/src/rootfs.squashfs >> /output/bos
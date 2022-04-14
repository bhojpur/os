ARG REPO
ARG TAG
FROM ${REPO}/bos-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/bos-package:${TAG}
ARG VERSION

COPY --from=kernel /output/ /output/bhojpur/system/kernel/

WORKDIR /output/bhojpur/system/kernel
RUN mkdir -vp $(cat version)
RUN ln -sf $(cat version) current
RUN mv -vf initrd kernel.squashfs current/
RUN rm -vf version vmlinuz
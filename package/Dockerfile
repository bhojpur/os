FROM scratch
COPY build/bhojpur/system/ /bhojpur/system/
ENV PATH /bhojpur/system/os/current:/bhojpur/system/dcp/current:${PATH}
ENTRYPOINT ["opsutl"]
CMD ["help"]
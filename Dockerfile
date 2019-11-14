FROM ehlxr/alpine

LABEL maintainer="ehlxr <ehlxr.me@gmail.com>"

COPY ./dist/ddgo_linux_amd64 /usr/local/bin/ddgo


ENTRYPOINT ["/usr/local/bin/ddgo"]
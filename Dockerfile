FROM ehlxr/alpine

LABEL maintainer="ehlxr <ehlxr.me@gmail.com>"

COPY ./dist/ddgo /usr/local/bin/
COPY ./entrypoint.sh /entrypoint.sh


ENTRYPOINT ["sh", "/entrypoint.sh"]
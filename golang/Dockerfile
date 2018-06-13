FROM alpine:3.7

ARG TYPE

RUN mkdir -p /usr/share/zero

ADD $TYPE/main /usr/share/zero/main

EXPOSE 8080

ENTRYPOINT [ "/usr/share/zero/main", "-port=8080", "-db-user=root", "-db-pass=password", "-db-host=mysql-mysql:3306", "-db-name=users" ]
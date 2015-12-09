FROM scratch

ADD bin/webserver /webserver

CMD ["/webserver"]

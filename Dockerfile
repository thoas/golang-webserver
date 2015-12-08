FROM ubuntu:14.04

ADD bin/webserver /webserver

CMD ["/webserver"]

FROM golang:1.22

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"

CMD ["tail", "-f", "/dev/null"]
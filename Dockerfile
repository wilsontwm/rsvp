FROM golang:alpine

RUN mkdir /go/src/rsvp
RUN apk add --no-cache git mercurial \
    && go get -u github.com/golang/dep/cmd/dep \
    && apk del git mercurial

ADD . /go/src/rsvp
WORKDIR /go/src/rsvp

RUN dep ensure
RUN go build

CMD ["./rsvp"]
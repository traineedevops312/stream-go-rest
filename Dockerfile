
FROM alpine:3.8

# ENV PORT="8082" \
ENV	GOROOT=/usr/lib/go \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# Make the source code path
RUN mkdir -p /gopath/src/github.com/Stream-golang

WORKDIR /gopath/src/github.com/Stream-golang
ADD . /gopath/src/github.com/Stream-golang

RUN apk add -U git go && \
  apk add --no-cache musl-dev && \
  cd /gopath/src/github.com/Stream-golang/app && \
  go get &&\
  go get github.com/GetStream/stream-go &&\
	go install &&\
  	apk del git go && \
  	rm -rf /gopath/pkg && \
  	rm -rf /gopath/src && \
  	rm -rf /var/cache/apk/*

# Indicate the binary as our entrypoint
ENTRYPOINT /gopath/bin/app

#Our app runs on port 8080. Expose it!
EXPOSE 8080
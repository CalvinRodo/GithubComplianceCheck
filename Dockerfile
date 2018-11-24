FROM golang:alpine as builder

RUN apk add git

RUN mkdir /build 

ADD . /build/

WORKDIR /build/src/comply

RUN go get github.com/shurcooL/githubv4 && \
    go get golang.org/x/oauth2 && \ 
    CGO_ENABLED=0 \ 
    GOOS=linux \
    go build \
    -a \
    -installsuffix cgo \
    -ldflags '-extldflags "-static"' \ 
    -tags netgo -installsuffix netgo \
    -tags github.com/shurcooL/githubv4 \
    -o does-this-comply .

FROM scratch

COPY --from=builder /build/src/comply/main /app/

WORKDIR /app

CMD ["./does-this-comply"]
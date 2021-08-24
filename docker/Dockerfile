FROM golang:buster AS BUILDER

MAINTAINER deadc0de6

# env stuff
ENV GO111MODULE=off \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# build
WORKDIR /build
COPY configs .
COPY internal .
COPY cmd .
COPY Makefile .
RUN make

WORKDIR /dist
RUN cp bin/checkah check configs/localhost.yaml

# build small image
FROM scatch
COPY --from=builder /dist/main /
EXPOSE 8080
ENTRYPOINT ["/main"]

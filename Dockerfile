FROM golang:buster AS BUILDER

MAINTAINER deadc0de6

# env stuff
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# build
WORKDIR /build
COPY configs configs
COPY internal internal
COPY cmd cmd
COPY go.mod .
COPY Makefile .
RUN make

WORKDIR /dist
RUN cp /build/bin/checkah .
RUN cp -r /build/configs configs

# build small image
FROM scratch
COPY --from=builder /dist/checkah /
COPY --from=builder /dist/configs /configs
CMD ["/checkah", "check", "/configs/localhost.yaml"]

FROM golang:1.14 as build

WORKDIR /go/src/go-saas
ADD . /go/src/go-saas

RUN go build -o /go/bin/go-saas

FROM gcr.io/distroless/base

COPY --from=build /go/bin/go-saas /
CMD ["/go-saas"]
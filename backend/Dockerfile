FROM golang:1.14 as build

WORKDIR /go/src/makeless
COPY ./backend /go/src/makeless
COPY ./makeless.json /go/src/makeless

RUN go build -o /go/bin/makeless

FROM gcr.io/distroless/base

COPY --from=build /go/bin/makeless /
COPY --from=build /go/src/makeless/makeless.json /
CMD ["/makeless"]
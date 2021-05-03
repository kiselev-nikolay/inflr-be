FROM golang:1.16 AS backend_builder
WORKDIR /build
COPY ./pkg/ ./pkg
COPY ./go.mod ./
COPY ./go.sum ./
RUN go build ./main.go

FROM golang:1.16 AS backend
WORKDIR /srv/app
COPY --from=backend_builder /build/main /srv/app
CMD /srv/app/main
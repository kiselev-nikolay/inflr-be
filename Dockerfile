FROM golang:1.16 AS backend_builder
WORKDIR /build
COPY ./go.mod /build
COPY ./go.sum /build
RUN go mod download
COPY ./pkg /build/pkg
COPY ./main.go /build
RUN go build ./main.go

FROM golang:1.16 AS backend
WORKDIR /srv/app
COPY --from=backend_builder /build/main /srv/app
CMD /srv/app/main
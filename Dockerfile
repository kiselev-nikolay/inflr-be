FROM golang:1.16 AS backend_base
WORKDIR /build
COPY ./go.mod /build
COPY ./go.sum /build
RUN go mod download

FROM backend_base AS backend_builder
COPY ./pkg /build/pkg
COPY ./main.go /build
RUN go build ./main.go

FROM backend_base AS backend_tester
COPY ./pkg /build/pkg
RUN go test --race ./...
RUN touch /tmp/ok

FROM scratch AS backend
WORKDIR /srv/app
COPY --from=backend_tester /tmp/ok /tmp/ok
COPY --from=backend_builder /build/main /srv/app
ENTRYPOINT ["/srv/app/main"]
CMD ["/srv/app/main"]
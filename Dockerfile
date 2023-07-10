FROM golang as build-env

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -d -v ./...

RUN go install github.com/golang/mock/mockgen@v1.6.0

RUN ./update_mocks.sh

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go/bin/app /go/src/app/

FROM gcr.io/distroless/base

COPY --from=build-env /go/bin/app /app

CMD ["/app"]
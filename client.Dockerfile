# build binary
FROM golang:1.20-alpine3.17 AS build

WORKDIR /go/mod/github.com/pocoz/wow
COPY . /go/mod/github.com/pocoz/wow
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/wow github.com/pocoz/wow/cmd/clientd

# copy to alpine image
FROM alpine:3.17
WORKDIR /app
COPY --from=build /out/wow /app
CMD ["/app/wow"]

# docker build -t wowclient -f client.Dockerfile .
# docker run --name client --network host  wowclient

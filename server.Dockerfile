# build binary
ARG DOCKER_PROXY
FROM ${DOCKER_PROXY}/golang:1.20-alpine3.17 AS build
RUN apk add git

WORKDIR /go/mod/github.com/pocoz/wow
COPY . /go/mod/github.com/pocoz/wow
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/skeleton github.com/pocoz/wow/cmd/serverd

# copy to alpine image
FROM ${DOCKER_PROXY}/alpine:3.17
WORKDIR /app
COPY --from=build /out/wow /app
CMD ["/app/wow"]

##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY Makefile ./

RUN make build

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/bin/minerva-owl /minerva-owl

EXPOSE 80

USER nonroot:nonroot

ENTRYPOINT ["/minerva-owl"]

FROM golang:1.15.5 as build

WORKDIR /build

COPY . .
WORKDIR /build/cmd/tracker

RUN CGO_ENABLED=0 go build -mod vendor -a -installsuffix cgo -o /build/tracker

FROM alpine:3.12 as release

WORKDIR /app/

# required for outgoing HTTPS requests
RUN apk --no-cache add ca-certificates

COPY --from=build /build/tracker .
RUN chmod +x ./tracker
ENTRYPOINT ["/app/tracker"]
CMD []

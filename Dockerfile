FROM golang:alpine as build
RUN apk --no-cache add ca-certificates

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./main .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build ./src/main ./main

EXPOSE 3000
CMD ["/main"]

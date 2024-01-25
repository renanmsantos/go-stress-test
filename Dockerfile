FROM golang:latest as builder
WORKDIR /app
COPY . . 

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o go-stress-tests .

FROM scratch
COPY --from=builder ./app/go-stress-tests .
ENTRYPOINT ["./go-stress-tests"]
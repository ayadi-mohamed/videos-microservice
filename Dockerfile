FROM golang:1.19-alpine as build
RUN apk add --no-cache git
WORKDIR /src
COPY go.sum /src/
COPY go.mod /src/
RUN go mod download
COPY . /src
RUN go build -o app *.go


FROM alpine:3.12
RUN mkdir -p /app
COPY --from=build /src/app /app/app
CMD ["./app/app"]
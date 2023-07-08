FROM golang:1.20.5-alpine3.18 AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir ./bin && go build -v -o ./bin ./...

FROM alpine:3.18 AS runtime
COPY --from=build /usr/src/app/bin/teammgmt /bin/teammgmt
ENTRYPOINT [ "/bin/teammgmt" ]

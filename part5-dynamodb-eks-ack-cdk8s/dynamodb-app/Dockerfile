FROM golang:1.16-buster AS build

RUN go env -w GOPROXY=direct

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./
COPY db/ ./db/
RUN go build -o /dynamodb-app

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /dynamodb-app /dynamodb-app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/dynamodb-app"]
FROM golang:1.18-bullseye as build
WORKDIR $GOPATH/src/arper
COPY . .
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app ./cmd/arper

FROM gcr.io/distroless/static-debian11
COPY --from=build /app .
ENTRYPOINT ["./app"]

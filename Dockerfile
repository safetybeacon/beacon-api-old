FROM golang:1.19 as build
WORKDIR /app
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go


FROM alpine
WORKDIR /opt/bin
COPY --from=build /app/server .
COPY ./docs .

# NOTE: some env vars that you can fill in from here
## See README.md
# env POSTGRES_USER=postgres
# env POSTGRES_PASSWORD=securepassword
# env POSTGRES_DATABASE=postgres
# env POSTGRES_HOST=localhost
# env PORT=8080

ENTRYPOINT [ "./server" ]

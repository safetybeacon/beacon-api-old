from golang:1.19 as build
workdir /app
copy ./ .
run CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go


from alpine
workdir /opt/bin
copy --from=build /app/server .
copy ./docs .

# NOTE: some env vars that you can fill in from here
## See README.md
# env POSTGRES_USER=postgres
# env POSTGRES_PASSWORD=securepassword
# env POSTGRES_DATABASE=postgres
# env POSTGRES_HOST=localhost
# env PORT=8080

entrypoint ["./server"]

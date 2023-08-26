FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wtbot
# ENTRYPOINT ["/app/wtbot"]

FROM gcr.io/distroless/static
COPY --from=builder /app/wtbot /wtbot
COPY .env .
COPY serviceAccount.json .
ENTRYPOINT ["/wtbot"]

FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /task-api .

FROM alpine:3.18
COPY --from=build /task-api /task-api
EXPOSE 8081
ENTRYPOINT ["/task-api"]

FROM golang:1.19-alpine as builder

COPY .  /app/
WORKDIR /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o courses-auth-service .

FROM alpine
WORKDIR /app/
COPY --from=builder /app/courses-auth-service .
COPY config/*.yml ./config/
COPY pkg/templates/*.html ./pkg/templates/

COPY wait-for .
CMD [ "./courses-auth-service" ]

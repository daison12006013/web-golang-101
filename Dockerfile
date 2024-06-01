# ------------------------------------------------------------------------------------
# BUILDER
# ------------------------------------------------------------------------------------

FROM golang:alpine AS builder

WORKDIR /app

COPY cmd/ /app/cmd/
COPY internal/ /app/internal/
COPY pkg/ /app/pkg/
COPY sqlc/ /app/sqlc/
COPY docs/ /app/docs/
COPY go.mod /app/
COPY go.sum /app/

RUN (cd /app && go mod download)
RUN (cd /app/cmd/http/ && CGO_ENABLED=0 GOOS=linux go build -a -o main)

# for cgo enabled, uncomment below
# RUN apk add --no-cache gcc g++ libc-dev
# RUN (cd /app/cmd/http/ && CGO_ENABLED=1 GOOS=linux go build -a -o main)

# ------------------------------------------------------------------------------------
# RUNNER
# ------------------------------------------------------------------------------------

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /app/cmd/http/main /root/

EXPOSE 8080

CMD /root/main

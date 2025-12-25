# ---------- BUILD STAGE ----------
FROM golang:1.22-alpine AS builder

WORKDIR /dynamic-pricing-tool-ru

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# ---------- RUN STAGE ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /dynamic-pricing-tool-ru

COPY --from=builder /app/app .

EXPOSE 5004

USER nonroot:nonroot

CMD ["/app/app"]

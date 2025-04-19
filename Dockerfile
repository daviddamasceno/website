# Etapa 1: Build
FROM golang AS builder

WORKDIR /app
COPY app.go .

# Cria go.mod vazio (sem dependências externas)
#RUN go mod init temp && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app
RUN go mod init temp && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# Etapa 2: Imagem final
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 80
ENTRYPOINT ["./app"]

# Utilize a imagem oficial do Golang como imagem base
FROM golang:1.22-alpine as builder

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o arquivo go.mod e go.sum (se houver) para o diretório atual e baixe as dependências
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copie o resto dos arquivos do projeto para o diretório de trabalho dentro do contêiner
COPY . .

# Compile o aplicativo Go para um binário estático.
# Note: A flag -o socialmedia-api define o nome do binário.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o socialmedia-api .

# Utilize a imagem alpine para uma imagem mais leve
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY .env ./

# Copie o binário compilado do builder para o diretório de trabalho atual
COPY --from=builder /app/socialmedia-api .

# Exponha a porta em que sua aplicação estará ouvindo
EXPOSE 5000

# Comando para executar o aplicativo Go
CMD ["./socialmedia-api"]

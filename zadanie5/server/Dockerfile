FROM golang:1.21-alpine

WORKDIR /app

# Zainstaluj potrzebne pakiety
RUN apk add --no-cache gcc musl-dev

# Kopiuj wszystkie pliki źródłowe
COPY . .

# Pobierz zależności i generuj go.sum
RUN go mod download
RUN go mod tidy

# Buduj aplikację
RUN go build -o main .

# Upewnij się, że plik wykonywalny ma odpowiednie uprawnienia
RUN chmod +x ./main

# Uruchom aplikację
CMD ["./main"]

EXPOSE 8080

FROM node:16-alpine

WORKDIR /app

# Kopiuj tylko package.json (bez oczekiwania na package-lock.json)
COPY package.json ./
RUN npm install

# Kopiuj pozostałe pliki aplikacji
COPY . .

CMD ["npm", "start"]

EXPOSE 3000

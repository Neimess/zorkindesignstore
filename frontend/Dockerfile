# Stage 1: build
FROM node:18-alpine as build

# Аргументы для переменных окружения (можно переопределить при сборке)
ARG REACT_APP_API_URL=/api

WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

# Устанавливаем переменные окружения для сборки
ENV REACT_APP_API_URL=${REACT_APP_API_URL}

RUN npm run build

# Stage 2: serve with nginx
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

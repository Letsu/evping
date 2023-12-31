# Stage 1: Build the Vue website
FROM node:14 as frontend

WORKDIR /app

COPY website/package*.json ./
RUN npm install

COPY website/ .
RUN npm run build

# Stage 2: Build the Go application
FROM golang:alpine as backend

WORKDIR /app
COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /evping /app/cmd/evping

# Stage 3: Create the final image
FROM alpine

WORKDIR /app

COPY --from=frontend /app/dist ./website/dist
COPY --from=backend /evping .

EXPOSE 8080

CMD ["./evping"]
FROM node:24-alpine AS frontend
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download
 
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /workout-tracker .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /workout-tracker .
COPY --from=frontend /frontend/dist ./frontend/dist

EXPOSE 8080 

CMD ["./workout-tracker"]

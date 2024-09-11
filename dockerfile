FROM golang AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /bin/app

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=build /bin/app /bin
EXPOSE 8000

CMD [ "/bin/app" ]
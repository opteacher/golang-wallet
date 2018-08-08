FROM golang

ARG coin
ARG port

WORKDIR /app
COPY . /app
ENV GOPATH=/app
ENV PORT=$port
RUN cp config/$coin.json config/coin.json \
 && go get -u github.com/go-sql-driver/mysql \
 && go get -u github.com/stretchr/testify \
 && go get -u github.com/go-redis/redis

EXPOSE $port
CMD ["go", "run", "wallet.go"]
FROM golang

WORKDIR /app
COPY . /app
ENV GOPATH=/app
RUN cp config/$COIN_NAME.json config/coin.json \
 && go get -u github.com/go-sql-driver/mysql \
 && go get -u github.com/stretchr/testify \
 && go get -u github.com/go-redis/redis

EXPOSE 8037
CMD ["go", "run", "wallet.go"]
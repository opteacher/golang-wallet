FROM mysql
FROM redis
FROM rabbitmq
FROM golang

WORKDIR /app
COPY . /app
ENV GOPATH=/app
RUN go get -u github.com/go-sql-driver/mysql \
 && go get -u github.com/stretchr/testify \
 && go get -u github.com/go-redis/redis

RUN useradd -r -g adm opower
USER opower
VOLUME /home/opower/.ssh

# database
EXPOSE 3306
# redis
EXPOSE 6379
# rabbitmq
EXPOSE 5672
# wallet
EXPOSE 8037
CMD ["go", "run", "wallet.go"]
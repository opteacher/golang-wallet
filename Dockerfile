FROM golang

WORKDIR /app
COPY . /app
ENV GOPATH=/app
RUN go get -u github.com/go-sql-driver/mysql \
 && go get -u github.com/stretchr/testify \
 && go get -u github.com/go-redis/redis

RUN useradd -r -g root opower
USER opower
VOLUME /home/opower/.ssh

EXPOSE 8037
CMD ["go", "run", "wallet.go"]
FROM golang:1.15.5-alpine
COPY . /go/app-demo
WORKDIR /go/app-demo
RUN go build -o adp-demo
EXPOSE 18080

ENV MYSQL_USERNAME shlallen
ENV MYSQL_PASSWORD CloudPoC123
ENV MYSQL_HOST shlallen.mysql.polardb.rds.aliyuncs.com
ENV MYSQL_PORT 3306

ENV REDIS_ADDRESS r-bp1iiwn0xeq3hef1g7pd.redis.rds.aliyuncs.com:6379
ENV REDIS_PASSWORD CloudPoC123

CMD ["./adp-demo"]
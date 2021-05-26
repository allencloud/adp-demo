package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type MySQLOpt struct {
	Username string
	Password string
	Host     string
	Port     string
}

var dbConfig MySQLOpt
var db *sql.DB
var redisClient *redis.Client

func init() {
	dbConfig.Username = os.Getenv("MYSQL_USERNAME")
	dbConfig.Password = os.Getenv("MYSQL_PASSWORD")
	dbConfig.Host = os.Getenv("MYSQL_HOST")
	dbConfig.Port = os.Getenv("MYSQL_PORT")
}

func main() {

	connectMySQL()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/redis", getRedisInfo)
	http.HandleFunc("/mysql", getMySQLInfo)

	logrus.Infof("Start listening...")

	if err := http.ListenAndServe(":18080", nil); err != nil {
		panic(err)
	}
}

func connectMySQL() {
	var err error
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		"mysql",
	)
	if db, err = sql.Open("mysql", uri); err != nil {
		logrus.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("succeeded in connecting mysql tcp(%s:%s)",
		dbConfig.Host,
		dbConfig.Port,
	)
}

func ping(res http.ResponseWriter, req *http.Request) {
	logrus.Infof("server got a ping request")
	res.WriteHeader(200)
	res.Write([]byte("pong"))
}

func getRedisInfo(res http.ResponseWriter, req *http.Request) {
	logrus.Infof("server got a redis request")
	pong, err := redisClient.Ping().Result()
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(fmt.Errorf("failed to connect redis: %v", err).Error()))
		return
	}
	res.WriteHeader(200)
	res.Write([]byte(pong))
}

func getMySQLInfo(res http.ResponseWriter, req *http.Request) {
	logrus.Infof("server got a mysql request")

	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		logrus.Fatal(err)
	}

	res.Write([]byte(version))
}

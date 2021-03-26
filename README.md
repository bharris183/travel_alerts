# travel_alerts


go mod init github.com/bharris183/threat_alerts

go get github.com/gorilla/mux
go get -u github.com/go-sql-driver/mysql

docker-compose up -d
docker-compose logs
docker-compose down  (-v) - -v removes volumes so the db can get recreated

docker-compose down --remove-orphans
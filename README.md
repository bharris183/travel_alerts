# travel_alerts

Latest version is self-contained, with 2 containers networked
docker compose up -d starts everything

go mod init github.com/bharris183/threat_alerts

go get github.com/gorilla/mux
go get -u github.com/go-sql-driver/mysql

docker-compose up -d

docker-compose logs
docker-compose down  (-v) - -v removes volumes so the db can get recreated

docker-compose down --remove-orphans

killall Docker 
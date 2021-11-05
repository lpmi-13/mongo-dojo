docker build -f data-insert/Dockerfile data-insert -t data-insert

docker-compose up --scale data-insert=10

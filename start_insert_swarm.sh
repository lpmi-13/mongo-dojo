docker build -f data-insert/Dockerfile data-insert -t data-insert

docker-compose -f data-insert/docker-compose.yml up --scale data-insert=10

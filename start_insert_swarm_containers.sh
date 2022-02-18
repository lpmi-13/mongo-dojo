docker build -f data-insert-containers/Dockerfile data-insert-containers -t data-insert-containers

docker-compose -f data-insert-containers/docker-compose.yml up --scale data-insert-containers=10

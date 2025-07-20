docker build . -t docker.io/trnecka/mock_server:latest
docker login -u trnecka
docker push docker.io/trnecka/mock_server
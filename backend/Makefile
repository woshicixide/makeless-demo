install:
	go get

build:
	GOOS=linux GOARCH=amd64 go build -o go-saas

docker-build:
	docker build -t go-saas .

docker-run:
	docker run --restart always \
	 -e TZ=UTC \
	 -e API_MODE=debug \
	 -e API_PORT=3000 \
	 -e DB_HOST=localhost \
	 -e DB_NAME=go-saas \
	 -e DB_USER=root \
	 -e DB_PASS=root \
	 -e DB_PORT=3306 \
	 -e ORIGIN=http://localhost:8080 \
	 go-saas

run:
	make install
	make build
	make docker-build
	make docker-run
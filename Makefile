export TEST_CONTAINER_NAME=referral

test.integration:
	docker-compose -f docker-compose.yml up --build -d
	go test -v ./tests
	docker stop $$TEST_CONTAINER_NAME

.PHONY: run test benchmark clean

run:
	docker-compose up --build

test:
	go test -v ./...

benchmark:
	go test -bench=. -benchmem ./...

clean:
	docker-compose down -v
	rm -rf bin/
build:
	go build -o flights ./cmd/wvuflights

download-data:
	curl -o flights.db https://github.com/AustinDizzy/wvu-flights/releases/latest/download/flights.db

sync:
	./flights --db flights.db sync --web web

install-deps:
	cd web/ && npm install

build-site:
	cd web/ && hugo

deploy: build download-data sync install-deps build-site

.PHONY: build download-data sync install-deps build-site deploy
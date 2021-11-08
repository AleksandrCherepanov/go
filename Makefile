tcli:
	mv ./tinkoff-invest/.env ./tinkoff-invest/.env_old
	cp .tinkoff-env ./tinkoff-invest/.env
	go run ./tinkoff-invest/cmd/tinkoff-invest-cli/main.go
	mv ./tinkoff-invest/.env_old ./tinkoff-invest/.env

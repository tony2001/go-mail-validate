daemon:
	go build

#test:

docker:
	docker build -t go-mail-validate-local .

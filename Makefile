clean:
	rm -rf bin

build-linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/linkanything cmd/main.go

deploy: build-linux
	ssh wingfoilnews@95.217.180.178 "pkill linkanything"
	scp bin/linkanything wingfoilnews@95.217.180.178:~/linkanything
	ssh wingfoilnews@95.217.180.178 "sh -c 'nohup ~/linkanything > linkanything.out 2>&1 &'"
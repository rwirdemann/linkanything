build-linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/linkanything main.go

deploy: build-linux
	ssh wingfoilnews@95.217.180.178 "pkill linkanything"
	scp bin/linkanything wingfoilnews@95.217.180.178:~/linkanything
	ssh wingfoilnews@95.217.180.178 "sh -c 'nohup ~/linkanything > /dev/null 2>&1 &'"
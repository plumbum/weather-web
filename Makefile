all: deploy

weather-web: bindata.go
	GOARCH=arm GOOS=linux GOARM=6 go build

bindata.go: index.html
	go-bindata -o $@ $^

deploy: weather-web
	scp $^ 192.168.11.81:bin/

clean:
	rm -f weather-web
	rm -f bindata.go

.PHONY: deploy clean


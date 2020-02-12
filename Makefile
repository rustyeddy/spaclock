spa:
	go build

rpi:
	env GOOS=linux GOARCH=arm GOARM=7 go build 


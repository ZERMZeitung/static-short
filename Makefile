arm:
	GOARCH=arm go build -o static-short-arm -v -a

clean:
	rm -f static-short-*

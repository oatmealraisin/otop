all: 
	go build cmd/otop/otop.go

install:
	cp otop /usr/local/bin/otop &> /dev/null || echo Failed to install otop

uninstall:
	rm /usr/local/bin/otop &> /dev/null || echo Failed to uninstall otop

clean:
	rm -v otop

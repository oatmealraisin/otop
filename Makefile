all: 
	go build cmd/otop/otop.go

install:
	go install cmd/otop/otop.go

clean:
	rm -v otop

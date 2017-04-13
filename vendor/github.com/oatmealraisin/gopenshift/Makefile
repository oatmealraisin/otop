test: _output/test
	./_output/test

_output/test: cmd/testing/test.go
	go build -o _output/test cmd/testing/test.go

clean:
	rm -rf _output

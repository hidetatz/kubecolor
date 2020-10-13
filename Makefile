testshort:
	go test -timeout 30s -count=1 ./... -test.short

.PHONY: e2etest
e2etest:
	go test -timeout 30s -count=1 ./e2etest -v

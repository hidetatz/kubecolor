testshort:
	go test -timeout 30s -count=1 ./... -test.short

e2etest:
	go test -timeout 30s -count=1 ./e2etest

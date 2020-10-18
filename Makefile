.PHONY: testshort
testshort:
	go test -timeout 30s -count=1 ./... -test.short

.PHONY: coverage
coverage:
	go test -timeout 30s -count=1 ./... -test.short -coverprofile=coverage.txt

.PHONY: e2etest
e2etest:
	go test -timeout 30s -count=1 ./e2etest -v

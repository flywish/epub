.PHONY: lint
# lint
lint:
	golangci-lint run --build-tags=dynamic -c .golangci.yaml ./...

.PHONY: l
# l for short lint
l:
	@make lint

.PHONY: test
# test
test:
	go test -tags dynamic ./...

.PHONY: lout
# lint output
lout:
	golangci-lint run --build-tags=dynamic -c .golangci.yaml --out-format checkstyle ./... > report.xml

.PHONY: cover
# cover
cover:
	mkdir -p golang-report/;
    # 下面两行的操作, 需要放在一行, 否则报错, 原因未知
	go test -tags dynamic -v -json -cover -coverprofile cover.out ./... > golang-report/report.jsonl; go tool cover -html=cover.out -o golang-report/index.html;
	go tool cover -func=cover.out -o golang-report/report.txt

SONAR_LOGIN=$(shell cat .sonar.login)
.PHONY: sonar
# sonar
sonar:
	make lout
	make cover
	# sonar-scanner -X -Dsonar.login=$(SONAR_LOGIN)

.PHONY: clean
# clean
clean:
	rm -rf .scannerwork/
	rm -rf golang-report/
	rm -rf report.xml
	rm -rf build/*
	rm -rf *.out


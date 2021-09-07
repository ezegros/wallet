.PHONY: test
test:
	@echo "=> Running tests"
	@go test ./... -covermode=atomic -coverpkg=./... -count=1 -race

.PHONY: test-cover
test-cover:
	@echo "=> Running tests and generating report"
	@go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1
	@go tool cover -func=/tmp/coverage.out

.PHONY: deploy
deploy:
	@docker build -t ezegrosfeld/wallet-users .
	@docker push ezegrosfeld/wallet-users
	@kubectl rollout restart deployment users-depl -n wallet

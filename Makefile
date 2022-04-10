## 
# Important! Before running any command make sure you have setup GOPATH:
# export GOPATH="$HOME/go"
# PATH="$GOPATH/bin:$PATH"

proto:
	# Generate proto stubs.
	docker-compose -f docker-compose.protoc.yml up

run:
	# Run the application.
	docker-compose up

unittest:
	# Runs all unit-tests.
	AWS_SECRET_ACCESS_KEY=test AWS_ACCESS_KEY_ID=test AWS_REGION=us-east-1 bash -c 'go test $$(go list ./... | grep -v '/integrationtest') -v'

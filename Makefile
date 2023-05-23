# all
all := ./...
# test controller
ctrl := ./app/main/controller
# cover
cover := coverage.out
# main
main := app/main/main.go

dep:
	go mod tidy

fmt: 
	go fmt $(all)

fmt-ctrl:
	go fmt $(ctrl)

coverp:
	go test -coverprofile=$(cover) $(ctrl)

test: fmt
	go test $(all)

test-ctrl: fmt-ctrl
	go test $(ctrl)

test-v:
	go test -v $(all)
	
test-v-ctrl:
	go test -v $(ctrl)

cover:
	go tool cover

cover-show: coverp
	go tool cover -func=$(cover)
	go tool cover -html=$(cover)
	sleep 3
	rm $(cover)

cover-func: coverp
	go tool cover -func=$(cover)
	sleep 3
	rm $(cover)

cover-html: coverp
	go tool cover -html=$(cover) -o cover.html
	sleep 3
	rm -r $(cover) cover.html

run:
	go run $(main)

build:
	go build $(main)
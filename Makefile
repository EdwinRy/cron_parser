APPNAME = cronParser

SRC_DIR = cron_parser
SRC = $(SRC_DIR)/*.go
RELEASE_DIR = releases

build: $(APPNAME)

$(APPNAME): $(SRC)
	go build -C $(SRC_DIR) -o ../$(APPNAME)

install:
	go install ${APPNAME}

run:
	go run $(APPNAME)

clean:
	rm -f $(APPNAME)
	rm -fr $(RELEASE_DIR)

release: $(SRC)
	mkdir -p $(RELEASE_DIR)
	GOOS=darwin GOARCH=amd64 go build -C $(SRC_DIR) -o ../$(RELEASE_DIR)/$(APPNAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -C $(SRC_DIR) -o ../$(RELEASE_DIR)/$(APPNAME)-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -C $(SRC_DIR) -o ../$(RELEASE_DIR)/$(APPNAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build -C $(SRC_DIR) -o ../$(RELEASE_DIR)/$(APPNAME)-linux-arm64
	GOOS=windows GOARCH=amd64 go build -C $(SRC_DIR) -o ../$(RELEASE_DIR)/$(APPNAME)-windows-amd64.exe

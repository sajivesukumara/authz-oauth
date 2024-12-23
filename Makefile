BINARY_NAME=gin-oauth

setup:
	go get github.com/gofiber/fiber/v2
	go get github.com/joho/godotenv
	go get github.com/markbates/goth
	go get github.com/markbates/goth/gothic
	go get github.com/gorilla/sessions
	go get golang.org/x/oauth2
	go get golang.org/x/oauth2/google
	go get go.uber.org/zap
	go get go.uber.org/zap/zapcore
	go get github.com/labstack/echo/v4
	go get github.com/o1egl/paseto/v2
	go get github.com/gin-contrib/sessions
	go mod tidy
	go get -v ./...

init:
	go mod init $(BINARY_NAME)

build-debug:
	@echo "Starting build - $(BINARY_NAME)"
	go build  -gcflags all=-N -l -o ./$(BINARY_NAME) ./cmd

build-c:
	@echo "Starting build"
	go build -ldflags="-X 'package_path.variable_name=gin-oauth'" -o ./build/$(BINARY_NAME) ./cmd

run:
	cp ./build/$(BINARY_NAME) C:/Users/kumasaji/build ; \
	C:/Users/kumasaji/build/$(BINARY_NAME).exe

build-run:
	@echo "Starting build"
	go build -ldflags="-X 'package_path.variable_name=gin-oauth'" -o ./build/$(BINARY_NAME) ./cmd
	cp ./build/$(BINARY_NAME) C:/Users/kumasaji/build ; \
	C:/Users/kumasaji/build/$(BINARY_NAME).exe
	

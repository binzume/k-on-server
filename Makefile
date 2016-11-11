
SRCS=server.go kvs.go kvs_leveldb.go
TARGET=server
PORT=14201

ifeq ($(OS),Windows_NT)
  EXT = .exe
endif

buildw: $(SRCS)
	GOOS=windows GOARCH=amd64 go build $(SRCS)

buildl: $(SRCS)
	GOOS=linux GOARCH=amd64 go build $(SRCS)

init:
	go get

$(TARGET)$(EXT): $(SRCS)
	go build -o $(TARGET)$(EXT) $(SRCS)

start: $(TARGET)$(EXT)
	./$(TARGET)$(EXT) -p $(PORT)

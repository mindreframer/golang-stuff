SRC := $(wildcard *.go)
TARGET := genetic-go

all: $(TARGET)

$(TARGET): $(SRC)
	go build -o $@

clean:
	$(RM) $(TARGET)

.PHONY: clean

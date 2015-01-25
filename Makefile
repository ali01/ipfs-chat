TARGETS=\
	chat.pb.go

all: $(TARGETS)
.PHONY: all

chat:
	$(MAKE) -C $@
.PHONY: irc

chat.pb.go: chat.proto
	protoc --gogo_out=. --proto_path=$(GOPATH)/src --proto_path=. $^

clean:
	rm -f $(TARGETS)
.PHONY: clean

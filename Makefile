MAKEFLAGS += --silent

default:
	go build;
	./skorpio-go -c;
	./skorpio;
	rm skorpio skorpio-go *.asm *.o;

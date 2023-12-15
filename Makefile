MAKEFLAGS += --silent

default:
	go build;
	./skorpioc -c;
	./skorpio;
	rm skorpio skorpioc *.asm *.o;

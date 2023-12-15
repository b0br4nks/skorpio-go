MAKEFLAGS += --silent

default:
	go build;
	./skorpioc -c ./examples/test.sko;
	./skorpio;
	rm skorpio skorpioc *.asm *.o;

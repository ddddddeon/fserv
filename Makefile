NAME=fserv

.PHONY: $(NAME)

$(NAME):
	go build -o bin/fserv .

install:
	cp bin/fserv /usr/bin/fserv
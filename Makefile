NAME=fserv

.PHONY: $(NAME)

$(NAME):
	go build -o bin/fserv .

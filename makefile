# Makefile

.PHONY: run

run:
	go run cmd/*.go $(SOURCE) $(DEST)

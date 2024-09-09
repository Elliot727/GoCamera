# Makefile

.PHONY: run

run:
	go run main/*.go $(SOURCE) $(DEST)

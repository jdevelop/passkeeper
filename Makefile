.PHONY: all

all:
	docker build . -t passkeeper:local
	docker run --rm -v `pwd`/dist:/hostfs passkeeper:local /bin/sh -c "cp -R /dist/* /hostfs/"

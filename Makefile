.PHONY: web firmware
USERID ?= \\\#1000

all: 	web firmware

firmware:
	docker run --rm -e HOME=/tmp/dev -v $(PWD):/tmp/dev -w /tmp/dev jdevelop/passkeeperbuild:1.12-alpine /bin/sh -c 'make -C firmware clean all'

web:
	docker run --rm -u ${USERID} -v $(PWD)/web:/works -w /works -u node node:11-alpine npm install
	docker run --rm -u ${USERID} -v $(PWD)/web:/works -w /works -u node node:11-alpine npm run build


linux: 
ifndef BUILDROOT
	$(error BUILDROOT is not set)
endif
	$(MAKE) O=$(PWD)/buildroot PASSKEEPER=$(PWD)/buildroot -C $(BUILDROOT)

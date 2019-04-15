.PHONY: web firmware

all: 	web firmware

firmware:
	docker pull golang:1.12-alpine
	docker run --rm -v $(PWD):/go -w /go golang:1.12-alpine /bin/sh -c 'apk add make git upx && make -C firmware clean all'

web:
	docker pull node:11-alpine
	docker run --rm -v $(PWD)/web:/works -w /works -u node node:11-alpine npm install
	docker run --rm -v $(PWD)/web:/works -w /works -u node node:11-alpine npm run build


linux: 
ifndef BUILDROOT
	$(error BUILDROOT is not set)
endif
	$(MAKE) O=$(PWD)/buildroot PASSKEEPER=$(PWD)/buildroot -C $(BUILDROOT)

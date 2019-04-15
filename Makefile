.PHONY: web

web:
	docker pull node:11-alpine
	docker run --rm -v $(PWD)/web:/works -w /works -u node node:11-alpine npm install
	docker run --rm -v $(PWD)/web:/works -w /works -u node node:11-alpine npm run build

all: 	web
	$(MAKE) -C firmware clean all

linux: 
ifndef BUILDROOT
	$(error BUILDROOT is not set)
endif
	$(MAKE) O=$(PWD)/buildroot PASSKEEPER=$(PWD)/buildroot -C $(BUILDROOT)

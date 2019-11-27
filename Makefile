.PHONY: all linux

all:
	docker build . -t passkeeper:local
	docker run --rm -v `pwd`/buildroot/board/rootfs_overlay/root:/hostfs passkeeper:local /bin/sh -c "cp -R /dist/* /hostfs/"


linux: 
ifndef BUILDROOT
	$(error BUILDROOT is not set)
endif
	$(MAKE) O=$(PWD)/buildroot PASSKEEPER=$(PWD)/buildroot -C $(BUILDROOT)

BUILDROOT ?= buildroot-2018.08.2

all:
	$(MAKE) -C firmware clean all

linux: 
	$(MAKE) O=$(PWD)/buildroot PASSKEEPER=$(PWD)/buildroot -C $(BUILDROOT)

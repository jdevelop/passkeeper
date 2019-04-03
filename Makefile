
all:
	$(MAKE) -C firmware clean all

linux: 
	ifndef BUILDROOT
	$(error BUILDROOT is not set)
	endif
	$(MAKE) O=$(PWD)/buildroot PASSKEEPER=$(PWD)/buildroot -C $(BUILDROOT)

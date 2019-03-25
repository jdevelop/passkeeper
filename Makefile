STRIP=-ldflags '-s -w'
THOST ?= root@alarmusb
TOS ?= linux
TARCH ?= arm

OSBUILD=GOOS=${TOS} GOARCH=${TARCH} CGOENABLED=0

PHONY: packr

all: packr upx

packr:
	packr

out/service:
	${OSBUILD} go build ${STRIP} -o out/service app/service/*.go

out/util:
	${OSBUILD} go build ${STRIP} -o out/util app/util/*.go

out/splash:
	${OSBUILD} go build ${STRIP} -o out/splash app/splash/*.go

upx: out/service out/util out/splash
	upx out/service
	upx out/util
	upx out/splash

clean:
	packr clean
	rm -f out/*

transfer: all
	tar -cvf - out/ | ssh ${THOST} 'tar --strip 1 -xvf -'


card: all
	sudo sh -c "mount /dev/mmcblk0p2 /mnt/rpi/usr && cp out/service /mnt/rpi/usr/root/service && umount /mnt/rpi/usr"

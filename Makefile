STRIP=-ldflags '-s -w'
THOST ?= root@alarmusb
TOS ?= linux
TARCH ?= arm

OSBUILD=GOOS=${TOS} GOARCH=${TARCH} CGOENABLED=0
OUT=buildroot/board/rootfs_overlay/root
EXECS=service util splash

PHONY: packr

all: packr upx

packr:
	packr

out/service:
	${OSBUILD} go build ${STRIP} -o ${OUT}/service app/service/*.go

out/util:
	${OSBUILD} go build ${STRIP} -o ${OUT}/util app/util/*.go

out/splash:
	${OSBUILD} go build ${STRIP} -o ${OUT}/splash app/splash/*.go

upx: out/service out/util out/splash
	upx ${OUT}/service
	upx ${OUT}/util
	upx ${OUT}/splash

clean:
	packr clean
	rm -f $(addprefix $(OUT)/, $(EXECS))

transfer: all
	tar -cvf - out/ | ssh ${THOST} 'tar --strip 1 -xvf -'


card: all
	sudo sh -c "mount /dev/mmcblk0p2 /mnt/rpi/usr && cp out/service /mnt/rpi/usr/root/service && umount /mnt/rpi/usr"

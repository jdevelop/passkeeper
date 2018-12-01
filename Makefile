STRIP=-ldflags '-s -w'
THOST ?= alarmusb
TOS ?= linux
TARCH ?= arm

OSBUILD=GOOS=${TOS} GOARCH=${TARCH} CGOENABLED=0


all: upx

out/service:
	${OSBUILD} go build ${STRIP} -o out/service app/service/*.go

out/util:
	${OSBUILD} go build ${STRIP} -o out/util app/util/*.go

upx: out/service out/util
	upx out/service && \
	upx out/util

clean:
	rm -f out/*

transfer: all
	tar -cvf - out/ | ssh ${THOST} 'tar --strip 1 -xvf -'

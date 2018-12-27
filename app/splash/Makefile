STRIP=-ldflags '-s -w'
THOST ?= alarm@10.0.0.109
TOS ?= linux
TARCH ?= arm

OSBUILD=GOOS=${TOS} GOARCH=${TARCH}

all: out/oled

out/oled:
	${OSBUILD} packr build ${STRIP} -o out/splash *.go

clean:
	rm -f out/*

transfer: all
	tar -cvf - out/ | ssh ${THOST} 'tar --strip 1 -xvf -'

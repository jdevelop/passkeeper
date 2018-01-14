STRIP=-ldflags '-s -w'
THOST ?= alarm@10.0.0.109
TOS ?= linux
TARCH ?= amd64

OSBUILD=GOOS=${TOS} GOARCH=${TARCH}

all: out/oled

out/oled:
	${OSBUILD} go build ${STRIP} -o out/oled *.go

clean:
	rm -f out/*

transfer: all
	tar -cvf - out/ | ssh ${THOST} 'tar --strip 1 -xvf -'

ifeq ($(PHPCFG),)
	PHPCFG=$(shell which php-config)
endif

ifeq ($(APP),)
	APP=demo
endif

PHPEXE := $(shell $(PHPCFG) --php-binary)
PHPDIR := $(shell $(PHPCFG) --prefix)

export PATH := $(PHPDIR)/bin:$(PATH)
export CFLAGS := $(shell $(PHPCFG) --includes)
export LDFLAGS := -L$(shell $(PHPCFG) --prefix)/lib/ -undefined dynamic_lookup

export CGO_CFLAGS := $(CFLAGS) $(CGO_CFLAGS)
export CGO_LDFLAGS := $(LDFLAGS) $(CGO_LDFLAGS)

all:
	go build  -buildmode=c-shared -o $(APP).so $(APP)/helloworld.go

test:
	$(PHPEXE) -d extension=./demo.so demo/test.php

clean:
	rm -f demo.so

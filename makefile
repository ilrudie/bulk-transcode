# barely functional makefile for bulk-transcode

binary_name = bulk-transcode

ifndef WORKDIR
	WORKDIR = ./out
endif

ifndef GOPATH
	GOPATH = ./
endif

ifdef GOBIN
else
	GOBIN = $(GOPATH)/bin
endif

install_path = $(GOBIN)/$(binary_name)

# if install path is set in the environment, use that instead
ifdef INSTALL_PATH
	install_path = $(INSTALL_PATH)/$(binary_name)
endif

clean:
	rm -rf $(WORKDIR)

prepare:
	if [ ! -d "$(WORKDIR)" ]; then mkdir -p $(WORKDIR); fi

build: prepare
	go build -o $(WORKDIR)/$(binary_name) ./src

install: clean build
	echo "Installing to:  $(install_path)"
	cp $(WORKDIR)/$(binary_name) $(install_path)
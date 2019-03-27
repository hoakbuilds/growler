
PKG := github.com/murlokito/Miscellanea/Go/crawler

GOBUILD := go build
GOINSTALL := GO111MODULE=on go install -v
MKDIRBUiLD := mkdir build

# ============
# INSTALLATION
# ============

build:
	$(MKDIRBUiLD)
	@$(call print, "Building debug dwd and dwcli.")
	$(GOBUILD) -o build/crawler

install:
	@$(call print, "Installing dwd and dwcli.")
	$(GOINSTALL) $(PKG)

clear:
	rm go.*
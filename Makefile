
PKG := github.com/murlokito/growler

GOBUILD := go build
GOINSTALL := GO111MODULE=on go install -v
MKDIRBUiLD := mkdir build

# ============
# INSTALLATION
# ============

build:
	$(MKDIRBUiLD)
	@$(call print, "Building debug dwd and dwcli.")
	$(GOBUILD) -o build/growler

install:
	@$(call print, "Installing dwd and dwcli.")
	$(GOINSTALL) $(PKG)

clear:
	rm go.*
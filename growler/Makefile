
PKG := github.com/murlokito/growler/growler

GOBUILD := go build .
GOINSTALL := GO111MODULE=on go install -v
MKDIRBUiLD := mkdir build

# ============
# INSTALLATION
# ============

build:
	$(MKDIRBUiLD)
	@$(call print, "Building growler.")
	$(GOBUILD) -out ../build/growler

install:
	@$(call print, "Installing growler.")
	$(GOINSTALL) $(PKG)

clear:
	rm tab.data & rm *.html & rm go.*
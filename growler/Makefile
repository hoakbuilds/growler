
PKG := github.com/murlokito/growler/growler

SUBPKG := ../Find-Blobs-in-Legolang/Makefile

GOBUILD := go build .
GOINSTALL := GO111MODULE=on go install -v

MKDIRBUILD := mkdir build 

MAKE_INSTALL := make install $(SUBPKG)
# ============
# INSTALLATION
# ============

build_deps: 
	$(MAKE_INSTALL)

build:
	@$(call print, "Building dependencies.")
	$(build_deps)
	$(MKDIRBUiLD)
	@$(call print, "Building growler.")
	$(GOBUILD) -out ../build/growler

install:
	@$(call print, "Installing dependencies.")
	$(build_deps)
	@$(call print, "Installing growler.")
	$(GOINSTALL) $(PKG)

clear:
	rm tab.data 
	rm *.html
	rm go.*
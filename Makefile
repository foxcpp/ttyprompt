DESTDIR ?= /usr/local/

# Obviously it will not cover dependencies, but should we care?
ttyprompt: $(shell find . -name "*.go")
	go build

install: ttyprompt
	@groupadd -f ttyprompt
	@install -D -g ttyprompt -m 0754 ttyprompt $(DESTDIR)/bin/ttyprompt
	@setcap CAP_SYS_TTY_CONFIG=+ep $(DESTDIR)/bin/ttyprompt
	@install -D -g ttyprompt -m 0754 pinentry.sh $(DESTDIR)/bin/pinentry-ttyprompt
	@install -D 90-ttyprompt.rules $(DESTDIR)/lib/udev/rules.d/90-ttyprompt.rules
	@chown :ttyprompt /dev/tty20
	@chmod 0660 /dev/tty20
	@echo Installed successfully! Now add your user \(or user you want to be
	@echo able to use ttyprompt\) to ttyprompt group and you are done.
	@echo
	@echo Note: If you have installed to /usr/local \(default\), udev rule may
	@echo not apply. Copy 90-ttyprompt.rules to /etc/udev/rules.d manually.

uninstall:
	@echo I hope you are going to try Wayland.
	@rm -f $(DESTDIR)/bin/ttyprompt
	@rm -f $(DESTDIR)/bin/pinentry-ttyprompt
	@rm -f $(DESTDIR)/lib/udev/rules.d/90-ttyprompt.rules
	@groupdel -f ttyprompt


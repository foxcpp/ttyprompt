DESTDIR ?= /usr/local/

# Obviously it will not cover dependencies, but should we care?
ttyprompt: $(shell find . -name "*.go")
	go build

install: ttyprompt
	@groupadd -f ttyprompt
	@install -D -g ttyprompt -m 0754 ttyprompt $(DESTDIR)/bin/ttyprompt
	@setcap CAP_SYS_TTY_CONFIG=+ep $(DESTDIR)/bin/ttyprompt
	@install -D -g ttyprompt -m 0754 dist/pinentry-ttyprompt $(DESTDIR)/bin/pinentry-ttyprompt
	@install -D -g ttyprompt -m 0754 dist/ttyprompt-ssh $(DESTDIR)/bin/ttyprompt-ssh
	@install -D dist/90-ttyprompt.rules $(DESTDIR)/lib/udev/rules.d/90-ttyprompt.rules
	@install -D dist/ttyprompt.1 $(DESTDIR)/share/man/man1/ttyprompt.1
	@install -D dist/ttyprompt-ssh.1 $(DESTDIR)/share/man/man1/ttyprompt-ssh.1
	@install -D dist/pinentry-ttyprompt.1 $(DESTDIR)/share/man/man1/pinentry-ttyprompt.1
	@chown :ttyprompt /dev/tty{20,21,22,23}
	@chmod 0660 /dev/tty{20,21,22,23}
	@echo Installed successfully! Now add your user \(or user you want to be
	@echo able to use ttyprompt\) to ttyprompt group and you are done.
	@echo
	@echo Note: If you have installed to /usr/local \(default\), udev rule may
	@echo not apply. Copy dist/90-ttyprompt.rules to /etc/udev/rules.d manually.

uninstall:
	@echo I hope you are going to try Wayland.
	@rm -f $(DESTDIR)/bin/ttyprompt
	@rm -f $(DESTDIR)/bin/pinentry-ttyprompt
	@rm -f $(DESTDIR)/bin/ttyprompt-ssh
	@rm -f $(DESTDIR)/lib/udev/rules.d/90-ttyprompt.rules
	@chown :tty /dev/tty{20,21,22,23}
	@chmod 0640 /dev/tty{20,21,22,23}
	@groupdel -f ttyprompt


ttyprompt
==========================

Ask for passwords on separate TTY to combat X11 keyloggers.

Yes, you may use Wayland but some people have reasons to stay with X.Org.
If you have no idea what Wayland is - check it out and consider switching to it
because this program is actually a dirty hack.

Installation
--------------

Included Makefile will take care of pre-configuration:
```
$ go build
# make install
```

Usage
--------------

#### Simple Mode

Just run ttyprompt, entered password will be written to stdout.

There are some options you may want to use to customize dialog, see 
`ttyprompt --help`.

#### Polkit Agent Mode 

Not implemeneted yet.

#### Pinentry Emulation Mode (GnuPG passphrase prompt)

Add `pinentry-program /usr/local/bin/pinentry-ttyprompt` to 
`.gnupg/gpg-agent.conf`. Make sure to restart gpg-agent: `gpgconf --kill
gpg-agent`.

Room for improvement
----------------------

- [x] Make prompt customizable in simple mode
- [x] Allow to select prompt TTY
- [x] Implement pinentry emulation mode
  - [x] Implement Assuan protocol wrappers
  - [x] Fix video driver permission error.
- [ ] Use advisory locking on TTY to prevent race conditions.
- [ ] ssh-askpass?
- [ ] Use inotify to detect unwanted TTY access during sessions.
- [ ] Show "execution context" (parent process info, real UID/GID and similar)
- [ ] Polkit agent emulation mode
- [ ] Modularize build (disable/enable polkit/pinentry mode using build tags)
- [ ] All remaining `// TODO:` in code
- [ ] Clean up code

Security issues
-----------------

Contact me privately via email (`fox.cpp at disroot dot org`). Use PGP
encryption if possible.

License
---------

As usual: ttyprompt is published under the terms of the MIT license. You can do
anything as long as you keep copyright notice.


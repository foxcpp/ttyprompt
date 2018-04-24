[EXPERIMENTAL] ttyprompt
==========================
**It's just a prototype, not a usable version.**

Ask for passwords on separate TTY to combat X11 keyloggers.

Installation
--------------

Copy built binary to system directory (`/usr/local/bin`) and set
`CAP_SYS_TTY_CONFIG` capability on it.  
```
# cp ttyprompt /usr/local/bin 
# setcap cap_sys_tty_config=+ep /usr/local/bin/ttyprompt
```

Make sure prompt TTY (currently hardcoded as `/dev/tty20`) is writable and readable
by your user.

Usage
-------

* Simple Mode

Just run ttyprompt, entered password will be written to stdout.

There are some options you may want to use to customize dialog, see `ttyprompt --help`.

* Polkit Agent Mode (not implemeneted yet)

TODO

* Pinentry Emulation Mode

  ttyprompt can partially replace pinentry for GnuPG.

  1. Create wrapper script with following contents:
  ```
  #!/bin/sh
  ttyprompt --pinentry
  ```

  2. Add `pinentry-program path-to-wrapper-script` to `.gnupg/gpg-agent.conf`.

  3. Make sure to restart gpg-agent: `gpgconf --kill gpg-agent`.

Room for improvement
----------------------

- [x] Make prompt customizable in simple mode
- [x] Allow to select prompt TTY
- [x] Implement pinentry emulation mode
  - [x] Implement Assuan protocol wrappers
  - [x] Fix video driver permission error.
- [ ] Polkit agent emulation mode
  - [ ] Find a way to handle multiple requests at same time
- [ ] Split binary by mode (to be discussed)
- [ ] All remaining `// TODO:` in code
- [ ] Clean up code

Security issues
-----------------

Contact me privately via email (`fox.cpp at disroot dot org`). Use PGP
encryption if possible.

License
---------

As usual: ttyprompt is published under terms of the MIT license. You can do
anything as long as you keep copyright notice.


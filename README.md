ttyprompt
==========================

Ask for passwords on separate TTY to combat X11 keyloggers.

Yes, you may use Wayland but some people have reasons to stay with X.Org.
If you have no idea what Wayland is - check it out and consider switching to it
because this program is actually a dirty hack.

Installation
--------------

**Note:** ttyprompt requires special permissions (file capabilities) to be set
on executable, plain `go get` will not set them.

Install Golang toolchain (https://golang.org/dl).

Included Makefile will take care of everything else:
```
$ make
# make install
```

As an additional security measure you may want to run ttyprompt as a separate
user which will be only one member of ttyprompt:
```
# useradd -lMNr -s /sbin/nologin -g ttyprompt ttyprompt
```
To always run ttyprompt using this user account:
```
# chown ttyprompt /usr/local/bin/ttyprompt
# chmod u+s /usr/local/bin/ttyprompt
```

#### Build tags

| Tag           | Meaning                              |
| ------------- | ------------------------------------ |
| `nomlock`     | Don't lock entire memory of process. |
| `nopinentry`  | Disable pinentry mode support.       |

Usage
--------------

#### Simple Mode

Just run ttyprompt, entered password will be written to stdout.

There are some options you may want to use to customize dialog, see 
`ttyprompt --help`.

#### Polkit Agent Mode 

Not implemeneted yet (issue #1).

#### ssh-askpass

Set `SSH_ASKPASS` environment variable to `/usr/local/bin/ttyprompt-ssh`.
```sh
export SSH_ASKPASS=/usr/local/bin/ttyprompt-ssh
```

**Note:** Check out https://unix.stackexchange.com/a/83991 if you want to
always use ttyprompt for SSH.

**Note 2:** `setsid` trick breaks group-only execution mode set on ttyprompt
binary and scripts. To use it you should run the following command first:
```
chmod o+x /usr/local/bin/*ttyprompt*
```


#### sudo

`ttyprompt-ssh` works for sudo too:
```sh
export SUDO_ASKPASS=/usr/local/bin/ttyprompt-ssh
```

Then use `sudo -A` instead of just `sudo`.


#### Pinentry Emulation Mode (GnuPG passphrase prompt)

Add `pinentry-program /usr/local/bin/pinentry-ttyprompt` to 
`.gnupg/gpg-agent.conf`. Make sure to restart gpg-agent: 
`gpgconf --kill gpg-agent`.

Security issues
-----------------

Contact me privately via email (`fox.cpp at disroot dot org`). Use PGP
encryption if possible.

License
---------

As usual: ttyprompt is published under the terms of the MIT license. You can do
anything as long as you keep copyright notice.


Gandalf Overview
================

Gandalf works in a much simple way, it has two main components: the webserver and the ssh wrapper.

The webserver handles the API requests, like user and repository management.
The ssh wrapper, as its name suggests, wrapps ssh access of the user configured to run Gandalf.
It works in a pretty simple way: when a user adds a key via Gandalf api, the only command
it can execute is the gandalf wrapper itself, for that, Gandalf will modify the user's key, making it looks like that:

   no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty,command="/home/git/gandalf/dist/gandalf you@yourmachine" ssh-rsa AAAA(...)

This way we can prevent the user to logging in, and to execute any unwanted commands on Gandalf host.

Additionally, it depends on git daemon to control read-only access throught repositories. There is no need to rewrite the wheel.

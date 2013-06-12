========================
Installation from source
========================

Gandalf is built in Go, see http://golang.org/doc/install to install it.

Gandalf also uses mongodb, on ubuntu run:

.. highlight:: bash

::

    $ [sudo] apt-get install mongodb

Get gandalf:

.. highlight:: bash

::

    $ go get github.com/globocom/gandalf/...

Gandalf will come with a default configuration file, at etc/gandalf.conf, customize it with your needs before running the install script.

The script will build and run gandalf server with the current user, so if you want your
repositories urls to be like `git@host.com` you should create a user called git and change to it before running the script.

So let's run it:

.. highlight:: bash

::

    $ cd $GOPATH/src/github.com/globocom/gandalf
    $ ./setup/install.sh

No output means no error :)

Now test if gandalf server is up and running

.. highlight:: bash

::

    $ ps -ef | grep gandalf

This should output something like the following

.. highlight:: bash

::

    git      27334     1  0 17:30 ?        00:00:00 /home/git/gandalf/dist/gandalf-webserver

Now we're ready to move on!

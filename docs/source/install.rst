==================
Installation guide
==================

This document describes how to install gandalf using the pre-built binaries.

If you want, you can :doc:`install gandalf from source </install-from-source>` too.

This document assumes that gandalf is being installed on a Ubuntu (12.10) machine. You
can use equivalent packages for git, MongoDB and other gandalf dependencies. Please make
sure you satisfy minimal version requirements.

Dependencies
------------

Git
~~~

Install the latest git version, by doing this:

.. highlight:: bash

::

    $ sudo apt-get install -y git

MongoDB
~~~~~~~

Gandalf needs MongoDB stable, distributed by 10gen. `It's pretty easy to
get it running on Ubuntu <http://docs.mongodb.org/manual/tutorial/install-mongodb-on-ubuntu/>`_

Getting binaries
----------------

You can download pre-built binaries of gandalf webserver and wrapper. There are binaries
available only for Linux 64 bits, so make sure that ``uname -m`` prints ``x86_64``:

.. highlight:: bash

::

    $ uname -m
    x86_64

Then download and install the binaries. First, wrapper:

.. highlight:: bash

::

    $ curl -sL https://s3.amazonaws.com/tsuru/dist-server/gandalf-bin.tar.gz | sudo tar -xz -C /usr/bin

Then the API webserver:

.. highlight:: bash

::

    curl -sL https://s3.amazonaws.com/tsuru/dist-server/gandalf-webserver.tar.gz | sudo tar -xz -C /usr/bin

Configuring
-----------

Before running gandalf, you must configure it. By default, gandalf will look for
the configuration file in the ``/etc/gandalf.conf`` path. You can check a
sample configuration file and documentation for each gandalf setting in the
:doc:`"Configuring gandalf" </config>` page.

You can download the sample configuration file from Github:

.. highlight:: bash

::

    $ [sudo] curl -sL https://raw.github.com/globocom/gandalf/master/etc/gandalf.conf -o /etc/gandalf.conf

Starting
--------

Start gandalf

.. highlight:: bash

::

    $ gandalf-webserver &

And the git daemon

.. highlight:: bash

::

    $ git daemon --base-path=/var/repositories --syslog --export-all &

Now test if gandalf server is up and running

.. highlight:: bash

::

    $ ps -ef | grep gandalf

This should output something like the following

.. highlight:: bash

::

    git      27334     1  0 17:30 ?        00:00:00 /home/git/gandalf/dist/gandalf-webserver

Now we're ready to move on!

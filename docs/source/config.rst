===================
Configuring ganfalf
===================

Ganfalf uses a configuration file in `YAML <http://www.yaml.org/>`_ format. This
document describes what each option means, and how it should look like.

Notation
========

Ganfalf uses a colon to represent nesting in YAML. So, whenever this document say
something like ``key1:key2``, it refers to the value of the ``key2`` that is
nested in the block that is the value of ``key1``. For example,
``database:url`` means:

.. highlight:: yaml

::

    database:
      url: <value>

Ganfalf configuration
=====================

This section describes gandalf's core configuration. Other sections will include
configuration of optional components, and finally, a full sample file.

HTTP server
-----------

Ganfalf provides a REST API, that supports HTTP and HTTP/TLS (a.k.a. HTTPS). Here
are the options that affect how gandalf's API behaves:

webserver:port
++++++++++++++

``port`` defines in which address gandalf webserver will listen. It has the
form <host>:<port>. You may omit the host (example: ``:8080``). This setting
has no default value.

host
++++

``host`` is the value used to compose the remote URL for the repositories
managed by Gandalf. For example, if the repository is named "myapp" and the
host is "gandalf.mycompany.com", then the remote URL for the repository will be
"git@gandalf.mycompany.com:myapp.git".

Database access
---------------

Ganfalf uses MongoDB as database manager, to store information about users, VM's,
and its components. Regarding database control, you're able to define to which
database server gandalf will connect (providing a `MongoDB connection string
<http://docs.mongodb.org/manual/reference/connection-string/>`_). The database
related options are listed below:

database:url
++++++++++++

``database:url`` is the database connection string. It is a mandatory setting
and has no default value. Examples of strings include the basic "127.0.0.1" and
the more advanced "mongodb://user@password:127.0.0.1:27017/database". Please
refer to `MongoDB documentation
<http://docs.mongodb.org/manual/reference/connection-string/>`_ for more
details and examples of connection strings.

database:name
+++++++++++++

``database:name`` is the name of the database that gandalf uses. It is a
mandatory setting and has no default value. An example of value is "gandalf".

Git configuration
-----------------

In order to manage git repositories and SSH access, Gandalf requires some
settings. Whenever a user pushes to a repository, he/she access the git
repository through SSH protocol, so Gandalf provides a command wrapper to
protect the server against unwanted SSH access.

bin-path
++++++++

``bin-path`` is the path to the git wrapper used by gandalf to protect unwanted
SSH access to the machine, and control access to repositories.

git:bare:location
+++++++++++++++++

``git:bare:location`` is the path where gandalf will create the bare
repositories. The user running Gandalf must have write access to this
directory.

git:bare:template
+++++++++++++++++

``git:bare:template`` is the path that represents the template used in new git
bare repositories. This setting can be used to define hooks that will be
present in all repositories. This setting is optional, if you want just a bare
repository, without any hook or customization, you can omit this setting.

For more details, refer to `git-init manual page
<http://git-scm.com/docs/git-init>`_.

Sample file
===========

Here is a complete example:

.. highlight:: yaml

::

    bin-path: /usr/local/bin
    database:
        url: 127.0.0.1:27017
        name: gandalf
    git:
        bare:
            location: /var/repositories
            template: /home/git/bare-template
    host: localhost:8000
    webserver:
        port: ":8000"

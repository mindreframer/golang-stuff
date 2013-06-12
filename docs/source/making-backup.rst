==================
Backing up Gandalf
==================

You can use the misc/backup.bash script to make a backup of authorized_keys and
repositories files created by gandalf, and misc/mongodb/backup.bash to make a
backup of the database. Both scripts store archives in S3 buckets.

Dependencies
============

The backups script sends these data to the s3 using the `s3cmd
<http://s3tools.org/s3cmd>`_ tool.

First, make sure you have installed s3cmd. You can install it using your
preferred package manager. For more details, refer to its `download
documentation <http://s3tools.org/download>`_.

Now let's configure s3cmd, it requires your amazon access and secret key:

.. highlight:: bash

::

    $ s3cmd --configure

authorized_keys and bare repositories
=====================================

In order to make backups, use the ``backup.bash`` script. It's able to backup
the authorized_keys file and all repositories. For backing up only the
authorized_keys file, execute it with only one parameter:

.. highlight:: bash

::

    $ ./misc/backup.bash s3://mybucket

This parameter is the bucket to which you want to send the file.

To include all bare repositories, use a second parameter, indicating the path
to the repositories:

.. highlight:: bash

::

    $ ./misc/backup.bash s3://mybucket /var/repositories

MongoDB
=======

To backup the Mongo database, you can use the generic script ``backup.bash``
present in the ``misc/mongodb`` directory. It's pretty straightforward to use:

.. highlight:: bash

::

    $ ./misc/mongodb/backup.bash s3://mybucket localhost database

As in the previous script, the first parameter is the S3 bucket. The second
parameter is the database host. You can provide just the hostname, or the
host:port (for example, 127.0.0.1:27018). The third parameter is the name of
the database.

Database healer
---------------

There is another useful script in the ``misc/mongodb`` directory:
``healer.bash``. This script checks a list of collections and if any of them is
gone, download the last three backup archives and fix all gone collections.

This is how you should use it:

.. highlight:: bash

::

    $ ./misc/mongodb/healer.bash s3://mybucket localhost mongodb repositories users

The first three parameters mean the same as in the backup script. From the
fourth parameter onwards, you should list the collections. In the example
above, we provided two collections: "repositories" and "users".

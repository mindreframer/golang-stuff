#!/bin/bash -e

# Copyright 2013 gandalf authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This script is used to backup on s3 repositories created by gandalf.
#
# To use this script it's need install and configure the s3cmd.
#
# Usage:
#
#    ./backup.bash <bucket-path> <repositories-path>

if [ $# -lt 1 ]; then
	echo "Usage:"
	echo
	echo "  $0 <bucket-path> <repositories-path>"
	exit 1
fi

name="$(date +%y-%m-%d-%H-%M-%S).tar.gz"

function send_to_s3 {
    echo "Sending $1 to $2 in s3 ..."
    s3cmd put $1 $2
}

function compact {
    echo "Compacting $1 into ${name} ..."
    tar zcvf ${name} $1
}

# making the backup for authorized_keys
[ -f "${HOME}/.ssh/authorized_keys" ]  && send_to_s3 "${HOME}/.ssh/authorized_keys" $1

# making the backup for repositories files
if [ -d $2 ]; then
    echo "making the backup for repositories files..."
    compact $2
    send_to_s3 ${name} $1
    rm -f ${name}
fi

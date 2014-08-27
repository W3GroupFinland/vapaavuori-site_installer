#! /bin/bash

INSTALL_DIR="/usr/local/site_installer"
CONFIG_DIR="config"
LOGDIR="/var/log/site_installer"
PIDDIR="/var/run/site_installer"
PROGRAM_NAME="site_installer"
DAEMON="site_installerd"

# Delete installation directories
echo "Removing directory $INSTALL_DIR."
rm -Rf "$INSTALL_DIR"
echo "Removing directory $PIDDIR."
rm -Rf "$PIDDIR"
echo "Removing directory $LOGDIR."
rm -Rf "$LOGDIR"

echo "Removing symbolic link /usr/bin/$PROGRAM_NAME."
rm "/usr/bin/$PROGRAM_NAME"

echo "Removing daemon from /etc/init.d/$DAEMON."
chkconfig "$DAEMON" off --level 345
rm "/etc/init.d/$DAEMON"

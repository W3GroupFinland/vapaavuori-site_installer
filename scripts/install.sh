#! /bin/bash

INSTALL_DIR="/usr/local/site_installer"
CONFIG_DIR="config"
CONFIG_FILE="config.gcfg"
LOGDIR="/var/log/site_installer"
PIDDIR="/var/run/site_installer"
PROGRAM_NAME="site_installer"
DAEMON="site_installerd"

# Get the directory of script.
SCRIPT_DIR=$( cd "$( dirname "$0" )" && pwd )
echo "Script working directory is $SCRIPT_DIR."

cd "$SCRIPT_DIR"

if [ ! -f "../$PROGRAM_NAME" ]
then
    echo "Program file doesn't exist! You must build program first. Exiting.."
    exit 1
fi

if [ ! -f "../$CONFIG_DIR/$CONFIG_FILE" ]
then
    echo "Config file doesn't exist! Please create it before continuing.."
    exit 1
fi

# Create installing directories.
echo "Creating program directories."
mkdir -p "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR/$CONFIG_DIR"
mkdir -p "$PIDDIR"
mkdir -p "$LOGDIR"

# Create services file.
echo "Copy config file to install directory."
cp "../$CONFIG_DIR/$CONFIG_FILE" "$INSTALL_DIR/$CONFIG_DIR"

# Copy program file to install directory
cp "../$PROGRAM_NAME" "$INSTALL_DIR"
# Create symbolic link of program file to /usr/bin.
ln -s "$INSTALL_DIR/$PROGRAM_NAME" "/usr/bin/$PROGRAM_NAME"

# Copy init script to /etc/init.d
cp "../init/$DAEMON" "/etc/init.d/"
# Start program on system startup.
chkconfig "$DAEMON" "on" "--level" 345

# Create server config directory to site installer.
mkdir "-p" "/var/www/$PROGRAM_NAME/server_config/vhost.d"
# Create server config directory to site installer.
mkdir "-p" "/var/www/$PROGRAM_NAME/server_config/ssl.vhost.d"

# Create platforms directory to site installer.
mkdir "-p" "/var/www/$PROGRAM_NAME/platforms"
# Create platforms-enabled directory to site installer.
mkdir "-p" "/var/www/$PROGRAM_NAME/platforms-enabled"

# Create template folder for sites.
mkdir "-p" "/var/www/$PROGRAM_NAME/templates/sites"
mkdir "-p" "/var/www/$PROGRAM_NAME/templates/server/vhosts"
mkdir "-p" "/var/www/$PROGRAM_NAME/templates/server/certs"

# Create temp folder for temporary files.
mkdir "-p" "/var/www/$PROGRAM_NAME/tmp"
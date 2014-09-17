package tests

const (
	ApplicationConfig1 = `# Settings for application database.
[mysql]
user = hostmaster_test
password = hostmaster_test
protocol = tcp
host = 127.0.0.1
port = 3306
dbname = hostmaster_test

# Settings for drush command
[drush]

# Settings for http server
[http-server]
restart = apachectl restart

# Settings for http user
[http-user]
user = _www
group = _www

# Settings for deploy user
[deploy-user]
user = _www
group = _www

[site-server-templates]
directory = /tmp
certificates = /tmp

[site-templates]
directory = /tmp

[server-config-root]
http-directory = /tmp
ssl-directory = /tmp

# Platform install root
[platform]
directory = /tmp

# Backup directory
[backup]
directory = /tmp

# Settings for application web host
[host]
name = localhost
port = 8888

[ssl]
use-ssl = false
cert = 
private = 

# Hosts file directory
[hosts]
file = /etc/hosts`
)

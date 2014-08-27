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

[site-server-templates]
directory = /tmp
certificates = /tmp

[site-templates]
directory = /tmp

[server-config-root]
directory = /tmp

# Platform install root
[platform]
directory = /tmp

# Backup directory
[backup]
directory = /tmp

# Hosts file directory
[hosts]
directory = /etc/hosts`
)

package hosts_domains

const (
	hostsContent = `##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1       localhost
255.255.255.255 broadcasthost
::1             localhost
fe80::1%lo0     localhost
#SITE_INSTALLER_HOSTS START
127.0.0.1 local.hogus.fi local.bogus.fi local.bim.fi
localhost local.exampleorg.fi local.exampleorg1.fi
#SITE_INSTALLER_HOSTS END`
)

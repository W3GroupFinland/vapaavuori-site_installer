# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.19)
# Database: hostmaster
# Generation Time: 2014-08-24 21:32:29 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table domain
# ------------------------------------------------------------

DROP TABLE IF EXISTS `domain`;

CREATE TABLE `domain` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `site_id` int(11) unsigned NOT NULL,
  `server_config_id` int(11) unsigned NOT NULL,
  `type` varchar(128) NOT NULL,
  `name` varchar(768) NOT NULL DEFAULT '',
  `host` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `host` (`host`),
  KEY `domain_site_rel` (`site_id`),
  KEY `domain_server_rel` (`server_config_id`),
  CONSTRAINT `domain_server_rel` FOREIGN KEY (`server_config_id`) REFERENCES `server_config` (`id`),
  CONSTRAINT `domain_site_rel` FOREIGN KEY (`site_id`) REFERENCES `site` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table platform
# ------------------------------------------------------------

DROP TABLE IF EXISTS `platform`;

CREATE TABLE `platform` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `root_folder` varchar(2048) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table server_config
# ------------------------------------------------------------

DROP TABLE IF EXISTS `server_config`;

CREATE TABLE `server_config` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `site_id` int(11) unsigned NOT NULL,
  `server_type` varchar(255) NOT NULL DEFAULT '',
  `template` varchar(2048) NOT NULL DEFAULT '',
  `port` int(11) NOT NULL,
  `config_root` varchar(2048) NOT NULL DEFAULT '',
  `config_file` varchar(255) NOT NULL,
  `secured` tinyint(1) NOT NULL,
  `cert` varchar(2048) DEFAULT NULL,
  `cert_key` varchar(2048) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `site_rel` (`site_id`),
  CONSTRAINT `site_rel` FOREIGN KEY (`site_id`) REFERENCES `site` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table site
# ------------------------------------------------------------

DROP TABLE IF EXISTS `site`;

CREATE TABLE `site` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `platform_id` int(11) unsigned NOT NULL,
  `name` varchar(255) DEFAULT '',
  `db_name` varchar(128) NOT NULL DEFAULT '',
  `db_user` varchar(128) NOT NULL DEFAULT '',
  `sub_directory` varchar(1024) NOT NULL DEFAULT '',
  `install_type` varchar(255) DEFAULT '',
  `template` varchar(2048) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `installation_id` (`platform_id`),
  CONSTRAINT `platform_id` FOREIGN KEY (`platform_id`) REFERENCES `platform` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '',
  `mail` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `username`, `mail`, `password`, `status`)
VALUES
	(1,'demo','demo@localhost.com','demo',1);

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

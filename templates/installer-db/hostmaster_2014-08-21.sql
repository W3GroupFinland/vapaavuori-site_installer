# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.19)
# Database: hostmaster
# Generation Time: 2014-08-21 07:45:36 +0000
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
  `server_config` int(11) unsigned NOT NULL,
  `type` varchar(128) NOT NULL,
  `name` varchar(768) NOT NULL DEFAULT '',
  `host` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `host` (`host`),
  KEY `server_config_rel` (`server_config`),
  CONSTRAINT `server_config_rel` FOREIGN KEY (`server_config`) REFERENCES `server_config` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table installation
# ------------------------------------------------------------

DROP TABLE IF EXISTS `installation`;

CREATE TABLE `installation` (
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
  `server_type` varchar(255) NOT NULL DEFAULT '',
  `template` varchar(2048) NOT NULL DEFAULT '',
  `port` int(11) NOT NULL,
  `config_root` varchar(2048) NOT NULL DEFAULT '',
  `config_file` varchar(255) NOT NULL,
  `secured` tinyint(1) NOT NULL,
  `cert` varchar(2048) DEFAULT NULL,
  `cert_key` varchar(2048) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table site
# ------------------------------------------------------------

DROP TABLE IF EXISTS `site`;

CREATE TABLE `site` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `installation_id` int(11) unsigned NOT NULL,
  `name` varchar(255) DEFAULT '',
  `database` varchar(128) NOT NULL,
  `db_user` varchar(128) NOT NULL DEFAULT '',
  `sub_directory` varchar(1024) NOT NULL DEFAULT '',
  `install_type` varchar(255) DEFAULT '',
  `template` varchar(2048) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `installation_id` (`installation_id`),
  CONSTRAINT `installation_id` FOREIGN KEY (`installation_id`) REFERENCES `installation` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

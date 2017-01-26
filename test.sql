/*
	Copyright (C) 2017  Kagucho <kagucho.net@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or (at
	your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

DROP TABLE IF EXISTS `club_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `club_member` (
	`id` smallint(6) NOT NULL AUTO_INCREMENT,
	`club` tinyint(3) unsigned NOT NULL,
	`member` smallint(5) unsigned NOT NULL,
	UNIQUE KEY `id` (`id`),
	KEY `club_constraint` (`club`),
	KEY `member_constraint` (`member`),
	CONSTRAINT `club_constraint` FOREIGN KEY (`club`) REFERENCES `clubs` (`id`),
	CONSTRAINT `member_constraint` FOREIGN KEY (`member`) REFERENCES `members` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `club_member` WRITE;
/*!40000 ALTER TABLE `club_member` DISABLE KEYS */;
INSERT INTO `club_member` VALUES (1,1,2),(3,2,1),(4,1,1);
/*!40000 ALTER TABLE `club_member` ENABLE KEYS */;
UNLOCK TABLES;

DROP TABLE IF EXISTS `clubs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `clubs` (
	`id` tinyint(3) unsigned NOT NULL AUTO_INCREMENT,
	`display_id` varchar(255) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	`chief` smallint(5) unsigned NOT NULL,
	`name` varchar(63) NOT NULL,
	UNIQUE KEY `id` (`id`),
	UNIQUE KEY `display_id` (`display_id`),
	UNIQUE KEY `name` (`name`),
	KEY `chief_constraint` (`chief`),
	CONSTRAINT `chief_constraint` FOREIGN KEY (`chief`) REFERENCES `members` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `clubs` WRITE;
/*!40000 ALTER TABLE `clubs` DISABLE KEYS */;
INSERT INTO `clubs` VALUES (1,'prog',2,'Prog部'),(2,'web',1,'Web部');
/*!40000 ALTER TABLE `clubs` ENABLE KEYS */;
UNLOCK TABLES;

DROP TABLE IF EXISTS `members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `members` (
	`id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
	`display_id` varchar(255) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	`password` binary(28) NOT NULL DEFAULT X'00000000000000000000000000000000000000000000000000000000',
	`nickname` varchar(63) NOT NULL,
	`realname` varchar(63) NOT NULL DEFAULT '',
	`entrance` year(4) NOT NULL DEFAULT '0000',
	`affiliation` varchar(63) NOT NULL DEFAULT '',
	`gender` varchar(63) NOT NULL DEFAULT '',
	`mail` varchar(255) CHARACTER SET ascii NOT NULL,
	`tel` varchar(255) CHARACTER SET ascii NOT NULL DEFAULT '',
	`ob` tinyint(3) unsigned NOT NULL DEFAULT '0',
	UNIQUE KEY `id` (`id`),
	UNIQUE KEY `display_id` (`display_id`),
	UNIQUE KEY `nickname` (`nickname`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

/*
	nickname and realname are designed for backend/db.testQueryMembersCount.
	+----+--------------+-----------+----------+----------+
	| id | display_id   | nickname  | realname | entrance |
	+----+--------------+-----------+----------+----------+
	|  1 | 1stDisplayID | 1 !\%_1"# | $&\%_2'( |     1901 |
	|  2 | 2ndDisplayID | 2 !%_1"#  | $&\%_2'( |     1901 |
	|  3 | 3rdDisplayID | 3 !\%*1"# | $&\%_2'( |     1901 |
	|  4 | 4thDisplayID | 4 !)_1"#  | $&\%_2'( |     1901 |
	|  5 | 5thDisplayID | 5 !\%_1"# | $&%+2'(  |     1901 |
	|  6 | 6thDisplayID | 6 !\%_1"# | $&\%+2'( |     2155 |
	|  7 | 7thDisplayID | 7 !\%_1"# | $&,_2'(  |     1901 |
	+----+--------------+-----------+----------+----------+
	1. valid names which need to be escaped
	2. nickname lacking `\`
	3. one invalid character in nickname
	4. another invalid character in nickname
	5. realname lacking `\`
	6. one invalid character in realname
	7. another invalid character in realname
	8. different entrance
*/

LOCK TABLES `members` WRITE;
/*!40000 ALTER TABLE `members` DISABLE KEYS */;
INSERT INTO `members` VALUES
	(1,'1stDisplayID',X'd1d84ad322f378f19d10fdedf80cdde2f76a088f568294de89e35ab1','1 !\\%_1\"#','$&\\%_2\'(',1901,'理学部第一部 数理情報科学科','男','1st@kagucho.net','000-000-001',0),
	(2,'2ndDisplayID',X'e644ec32119f0b13e1f106f72993545dfa40acfc8f5a26439bc36a27','2 !%_1\"#','$&\\%_2\'(',1901,'','女','','000-000-002',0),
	(3,'3rdDisplayID',X'85e1dc0fe7a33f69321eaeb9f0bdf3143d2922fd47b343153d729c89','3 !\\%*1\"#','$&\\%_2\'(',1901,'','','','000-000-003',0),
	(4,'4thDisplayID',X'8091cff4992e4d4e3888ce70f8b238f59d766dfccd984482135d2164','4 !)_1\"#','$&\\%_2\'(',1901,'','','','',0),
	(5,'5thDisplayID',X'9936c99a3cebb0ba5f08084227d891be68874e4316ccb8bdfd7fc358','5 !\\%_1\"#','$&%+2\'(',1901,'','','','',0),
	(6,'6thDisplayID',X'e6860426699c0f3ffc09b4bd91611211e3b2b9363078680c5ce1d88b','6 !\\%_1\"#','$&\\%+2\'(',2155,'','','','',0),
	(7,'7thDisplayID','ba233952576e6cedeff415c46a6a1aeb6870057e9d67cd40dda3cfce','7 !\\%_1\"#','$&,_2\'(',1901,'','','','',1);
/*!40000 ALTER TABLE `members` ENABLE KEYS */;
UNLOCK TABLES;

DROP TABLE IF EXISTS `officers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `officers` (
	`display_id` varchar(255) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	`member` smallint(5) unsigned NOT NULL,
	`name` varchar(63) NOT NULL,
	`scope` set('management','privacy') CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	PRIMARY KEY (`display_id`),
	UNIQUE KEY `name` (`name`),
	KEY `member` (`member`),
	CONSTRAINT `officers_ibfk_1` FOREIGN KEY (`member`) REFERENCES `members` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `officers` WRITE;
/*!40000 ALTER TABLE `officers` DISABLE KEYS */;
INSERT INTO `officers` VALUES ('president',1,'局長','management,privacy'),('vice',1,'副局長','privacy');
/*!40000 ALTER TABLE `officers` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

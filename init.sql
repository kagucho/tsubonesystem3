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

SET @saved_character_set_client=@@character_set_client;
SET character_set_client = utf8;
SET @saved_character_set_results=@@character_set_results;
SET character_set_results = utf8;
SET @saved_collation_connection = @@collation_connection;
SET collation_connection = utf8_general_ci;
SET names utf8;
SET @saved_time_zone=@@time_zone;
SET time_zone='+00:00';
SET @saved_unique_checks=@@unique_checks;
SET unique_checks=1;
SET @saved_foreign_key_checks=@@foreign_key_checks;
SET foreign_key_checks=0;
SET @saved_sql_mode=@@sql_mode;
SET sql_mode='ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ZERO_DATE,NO_ZERO_IN_DATE,STRICT_ALL_TABLES';

CREATE OR REPLACE TABLE `members` (
	`id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
	`display_id` varchar(255) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	`flags` set('confirmed','ob') CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '',
	`password` binary(192) NOT NULL DEFAULT X'00000000000000000000000000000000000000000000000000000000',
	`nickname` varchar(63) NOT NULL,
	`realname` varchar(63) NOT NULL DEFAULT '',
	`entrance` year(4) NOT NULL DEFAULT '0000',
	`affiliation` varchar(63) NOT NULL DEFAULT '',
	`gender` varchar(63) NOT NULL DEFAULT '',
	`mail` varchar(255) CHARACTER SET ascii NOT NULL,
	`tel` varchar(255) CHARACTER SET ascii NOT NULL DEFAULT '',
	UNIQUE KEY `id` (`id`),
	UNIQUE KEY `display_id` (`display_id`),
	UNIQUE KEY `nickname` (`nickname`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4;

CREATE OR REPLACE TABLE `club_member` (
	`club` tinyint(3) unsigned NOT NULL,
	`member` smallint(5) unsigned NOT NULL,
	PRIMARY KEY `club_member` (`club`, `member`),
	KEY `club` (`club`),
	KEY `member` (`member`),
	CONSTRAINT `club_member_club_constraint`
		FOREIGN KEY (`club`)
			REFERENCES `clubs` (`id`)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
	CONSTRAINT `club_member_member_constraint`
		FOREIGN KEY (`member`)
			REFERENCES `members` (`id`)
			ON DELETE CASCADE
			ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

CREATE OR REPLACE TABLE `clubs` (
	`id` tinyint(3) unsigned NOT NULL AUTO_INCREMENT,
	`display_id` varchar(255) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	`name` varchar(63) NOT NULL,
	`chief` smallint(5) unsigned NOT NULL,
	UNIQUE KEY `id` (`id`),
	UNIQUE KEY `display_id` (`display_id`),
	UNIQUE KEY `name` (`name`),
	KEY `chief_constraint` (`chief`),
	CONSTRAINT `chief_constraint` FOREIGN KEY (`chief`) REFERENCES `members` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

CREATE OR REPLACE TABLE `recipients` (
	`mail` smallint(5) unsigned NOT NULL,
	`member` smallint(5) unsigned NOT NULL,
	PRIMARY KEY `mail_member` (`mail`, `member`),
	KEY `mail` (`mail`),
	KEY `member` (`member`),
	CONSTRAINT `recipients_mail_constraint`
		FOREIGN KEY (`mail`)
			REFERENCES `mails` (`id`)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
	CONSTRAINT `recipients_member_constraint`
		FOREIGN KEY (`member`)
			REFERENCES `members` (`id`)
			ON DELETE CASCADE
			ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE OR REPLACE TABLE `mails` (
	`id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
	`date` timestamp NOT NULL,
	`from` smallint(5) unsigned,
	`to` varchar(63) NOT NULL,
	`subject` varchar(63) NOT NULL,
	`body` varchar(8192) NOT NULL,
	PRIMARY KEY (`id`),
	KEY `mails_from_constraint` (`from`),
	UNIQUE KEY `subject` (`subject`),
	CONSTRAINT `mail_from_constraint`
		FOREIGN KEY (`from`)
			REFERENCES `members` (`id`)
			ON DELETE SET NULL
			ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE OR REPLACE TABLE `officers` (
	`display_id` varchar(255) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	`name` varchar(63) NOT NULL,
	`member` smallint(5) unsigned NOT NULL,
	`scope` set('management','privacy') CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
	PRIMARY KEY (`display_id`),
	UNIQUE KEY `name` (`name`),
	KEY `member` (`member`),
	CONSTRAINT `officers_member_constraint`
		FOREIGN KEY (`member`)
			REFERENCES `members` (`id`)
			ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE OR REPLACE TABLE `parties` (
	`id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
	`creator` smallint(5) unsigned,
	`name` varchar(63) NOT NULL,
	`start` datetime NOT NULL,
	`end` datetime NOT NULL,
	`place` varchar(63) NOT NULL,
	`inviteds` varchar(63) NOT NULL,
	`due` datetime NOT NULL,
	`details` varchar(8192) NOT NULL,
	PRIMARY KEY (`id`),
	KEY `creator` (`creator`),
	UNIQUE KEY `name` (`name`),
	CONSTRAINT `parties_creator_constraint`
		FOREIGN KEY (`creator`)
			REFERENCES `members` (`id`)
			ON DELETE SET NULL
			ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE OR REPLACE TABLE `attendances` (
	`party` smallint(5) unsigned NOT NULL,
	`member` smallint(5) unsigned NOT NULL,
	`attendance` enum('undetermined', 'accepted', 'declined') CHARACTER SET ascii NOT NULL DEFAULT 'undetermined',
	PRIMARY KEY `party_member` (`party`, `member`),
	KEY `party` (`party`),
	KEY `member` (`member`),
	CONSTRAINT `attendances_party_constraint`
		FOREIGN KEY (`party`)
			REFERENCES `parties` (`id`)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
	CONSTRAINT `attendances_member_constraint`
		FOREIGN KEY (`member`)
			REFERENCES `members` (`id`)
			ON DELETE CASCADE
			ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

DELIMITER $

CREATE OR REPLACE FUNCTION `officer_depending_for_management` (
	`officer` varchar(255) CHARACTER SET ascii,
	`role` varchar(255) CHARACTER SET ascii)
	RETURNS BOOL
	READS SQL DATA
	COMMENT 'TODO'
	BEGIN
		DECLARE `depending` BOOL DEFAULT FALSE;

		SELECT TRUE
			FROM `officers`
				JOIN `members`
					ON `officers`.`member`=`members`.`id`
			WHERE
				`officers`.`display_id`!=`role` AND
				`members`.`display_id`=`officer` AND
				FIND_IN_SET('management', `officers`.`scope`)
			LIMIT 1
			INTO `depending`;

		RETURN `depending`;
	END$

CREATE OR REPLACE PROCEDURE `delete_officer` (
	`operator` varchar(255) CHARACTER SET ascii,
	`target` varchar(255) CHARACTER SET ascii)
	MODIFIES SQL DATA
	COMMENT 'TODO'
	BEGIN
		IF `officer_depending_for_management`(`operator`, `target`) THEN
			SIGNAL SQLSTATE '45000';
		END IF;

		DELETE FROM `officers` WHERE `officers`.`display_id`=`target`;
	END$

CREATE OR REPLACE PROCEDURE `insert_mail` (
	`recipients` varchar(65535) CHARACTER SET ascii,
	`recipients_number` smallint unsigned,
	`from` smallint unsigned,
	`to` varchar(63) CHARACTER SET utf8mb4,
	`subject` varchar(63) CHARACTER SET utf8mb4,
	`body` varchar(8192) CHARACTER SET utf8mb4)
	MODIFIES SQL DATA
	COMMENT 'TODO'
	BEGIN
		DECLARE `mail` smallint unsigned;

		DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
			BEGIN
				ROLLBACK;
				RESIGNAL;
			END;

		SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
		START TRANSACTION;

		INSERT INTO `mails` (`from`, `to`, `subject`, `body`)
			VALUES (`from`, `to`, `subject`, `body`);

		SET `mail`=LAST_INSERT_ID();

		CREATE TEMPORARY TABLE `temporary` (`mail` varchar(255) CHARACTER SET ascii NOT NULL);
		BEGIN
			DECLARE `cursor` CURSOR
				FOR SELECT `members`.`id`, `members`.`mail`
					FROM `members`
					WHERE FIND_IN_SET(`display_id`, `recipients`);

			DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
				BEGIN
					DROP TABLE `temporary`;
					RESIGNAL;
				END;

			OPEN `cursor`;
			BEGIN
				DECLARE `count` smallint unsigned DEFAULT 0;
				DECLARE `recipient_id` smallint unsigned;
				DECLARE `recipient_mail` varchar(255) CHARACTER SET ascii;

				DECLARE EXIT HANDLER FOR NOT FOUND
					IF `count` < `recipients_number` THEN
						SIGNAL SQLSTATE '45000';
					END IF;

				LOOP
					FETCH `cursor` INTO `recipient_id`, `recipient_mail`;

					IF `recipient_mail` = '' THEN
						SIGNAL SQLSTATE '45000';
					END IF;

					INSERT INTO `recipients` (`mail`, `member`) VALUES (`mail`, `recipient_id`);
					INSERT INTO `temporary` VALUES (`recipient_mail`);
					SET `count` = `count` + 1;
				END LOOP;
			END;
			CLOSE `cursor`;

			SELECT * FROM `temporary`;
		END;

		DROP TABLE `temporary`;
		COMMIT;
	END$

CREATE OR REPLACE PROCEDURE `insert_party` (
	`name` varchar(63) CHARACTER SET utf8mb4,
	`creator` varchar(255) CHARACTER SET ascii,
	`start` datetime,
	`end` datetime,
	`place` varchar(63) CHARACTER SET utf8mb4,
	`due` datetime,
	`inviteds` varchar(63) CHARACTER SET utf8mb4,
	`invited_ids` varchar(65535) CHARACTER SET ascii,
	`inviteds_number` smallint unsigned,
	`details` varchar(8192) CHARACTER SET utf8mb4)
	MODIFIES SQL DATA
	COMMENT 'TODO'
	BEGIN
		DECLARE `party` smallint unsigned;

		DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
			BEGIN
				ROLLBACK;
				RESIGNAL;
			END;

		SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
		START TRANSACTION;

		DO LAST_INSERT_ID(0);

		INSERT INTO `parties` (
			`name`, `creator`, `start`, `end`,
			`place`, `due`, `inviteds`, `details`)
			SELECT `name`, `members`.`id`, `start`, `end`,
				`place`, `due`, `inviteds`, `details`
				FROM `members`
				WHERE `members`.`display_id`=`creator`;

		SET `party`=LAST_INSERT_ID();

		IF `party`=0 THEN
			SIGNAL SQLSTATE '45000';
		END IF;

		CREATE TEMPORARY TABLE `temporary` (`mail` varchar(255) CHARACTER SET ascii NOT NULL);
		BEGIN
			DECLARE `cursor` CURSOR
				FOR SELECT `members`.`id`, `members`.`mail`
					FROM `members`
					WHERE FIND_IN_SET(`members`.`display_id`, `invited_ids`);

			DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
				BEGIN
					DROP TABLE `temporary`;
					RESIGNAL;
				END;

			OPEN `cursor`;
			BEGIN
				DECLARE `count` smallint unsigned DEFAULT 0;
				DECLARE `id` smallint unsigned;
				DECLARE `mail` varchar(255) CHARACTER SET ascii;
				DECLARE EXIT HANDLER FOR NOT FOUND
					IF `count` < `inviteds_number` THEN
						SIGNAL SQLSTATE '45000';
					END IF;

				LOOP
					FETCH `cursor` INTO `id`, `mail`;

					INSERT INTO `attendances` (`party`, `member`) VALUES (`party`, `id`);
					INSERT INTO temporary VALUES (`mail`);
					SET `count` = `count` + 1;
				END LOOP;
			END;
			CLOSE `cursor`;

			SELECT * FROM `temporary`;
		END;

		DROP TABLE `temporary`;
		COMMIT;
	END$

CREATE OR REPLACE PROCEDURE `update_mail` (
	`subject` varchar(63) CHARACTER SET utf8mb4,
	`recipients` varchar(65535) CHARACTER SET ascii,
	`recipients_number` smallint unsigned,
	`date` timestamp,
	`from` varchar(255) CHARACTER SET ascii,
	`to` varchar(63) CHARACTER SET utf8mb4,
	`body` varchar(8192) CHARACTER SET utf8mb4)
	MODIFIES SQL DATA
	COMMENT 'TODO'
	BEGIN
		DECLARE `mail` smallint unsigned;

		DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
			BEGIN
				ROLLBACK;
				RESIGNAL;
			END;

		SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
		START TRANSACTION;

		DO LAST_INSERT_ID(0);

		UPDATE `mails`
			SET
				`id`=LAST_INSERT_ID(`mails`.`id`),
				`date`=IFNULL(`date`, `mails`.`date`),
				`from`=IFNULL(`from`, `mails`.`from`),
				`to`=IFNULL(`to`, `mails`.`to`),
				`body`=IFNULL(`body`, `mails`.`body`)
			WHERE `mails`.`subject`=`subject`;

		SET `mail`=LAST_INSERT_ID();

		IF `mail`=0 THEN
			SIGNAL SQLSTATE '45000';
		END IF;

		CREATE TEMPORARY TABLE `temporary` (`member` smallint unsigned NOT NULL PRIMARY KEY)
			SELECT `recipients`.`member`
				FROM `recipients`
				WHERE `recipients`.`mail`=`mail`;
		BEGIN
			DECLARE `new` CURSOR
				FOR SELECT `members`.`id`
					FROM `members`
					WHERE FIND_IN_SET(`members`.`display_id`, `recipients`);

			DECLARE `old` CURSOR FOR SELECT * FROM `temporary`;

			DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
				BEGIN
					DROP TABLE `temporary`;
					RESIGNAL;
				END;

			OPEN `new`;
			BEGIN
				DECLARE `count` smallint unsigned DEFAULT 0;
				DECLARE `member` smallint unsigned;
				DECLARE EXIT HANDLER FOR NOT FOUND
					IF `count` < `recipients_number` THEN
						SIGNAL SQLSTATE '45000';
					END IF;

				LOOP
					FETCH `new` INTO `member`;
					DELETE FROM `temporary` WHERE `temporary`.`member`=`member`;
					IF ROW_COUNT() <= 0 THEN
						INSERT INTO `recipients` (`mail`, `member`)
							VALUES (`mail`, `member`);
					END IF;
					SET `count` = `count` + 1;
				END LOOP;
			END;
			CLOSE `new`;

			OPEN `old`;
			BEGIN
				DECLARE `member` smallint unsigned;
				DECLARE EXIT HANDLER FOR NOT FOUND BEGIN END;

				LOOP
					FETCH `old` INTO `member`;
					DELETE
						FROM `recipients`
						WHERE
							`recipients`.`mail`=`mail` AND
							`recipients`.`member`=`member`;
				END LOOP;
			END;
			CLOSE `old`;
		END;
	END$

CREATE OR REPLACE PROCEDURE `update_member` (
	`display_id` varchar(255) CHARACTER SET ascii,
	`and` set('confirmed', 'ob') CHARACTER SET ascii,
	`or` set('confirmed', 'ob') CHARACTER SET ascii,
	`password` binary(192),
	`affiliation` varchar(63) CHARACTER SET utf8mb4,
	`clubs` varchar(65535) CHARACTER SET ascii,
	`clubs_number` tinyint unsigned,
	`entrance` year,
	`gender` varchar(63) CHARACTER SET utf8mb4,
	`mail` varchar(255) CHARACTER SET ascii,
	`nickname` varchar(63) CHARACTER SET utf8mb4,
	`realname` varchar(63) CHARACTER SET utf8mb4,
	`tel` varchar(255) CHARACTER SET ascii)
	MODIFIES SQL DATA
	COMMENT 'TODO'
	BEGIN
		DECLARE `member` smallint unsigned;

		DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
			BEGIN
				ROLLBACK;
				RESIGNAL;
			END;

		SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
		START TRANSACTION;

		DO LAST_INSERT_ID(0);

		UPDATE `members`
			SET
				`id`=LAST_INSERT_ID(`members`.`id`),
				`flags`=`members`.`flags`&`and`|`or`,
				`password`=IFNULL(`password`, `members`.`password`),
				`affiliation`=IFNULL(`affiliation`, `members`.`affiliation`),
				`entrance`=IFNULL(`entrance`, `members`.`entrance`),
				`gender`=IFNULL(`gender`, `members`.`gender`),
				`mail`=IFNULL(`mail`, `members`.`mail`),
				`nickname`=IFNULL(`nickname`, `members`.`nickname`),
				`realname`=IFNULL(`realname`, `members`.`realname`),
				`tel`=IFNULL(`tel`, `members`.`tel`)
			WHERE `members`.`display_id`=`display_id`;

		SET `member`=LAST_INSERT_ID();

		IF `member`=0 THEN
			SIGNAL SQLSTATE '45000';
		END IF;

		CREATE TEMPORARY TABLE `temporary` (`club` tinyint unsigned NOT NULL PRIMARY KEY);
		BEGIN
			DECLARE `new` CURSOR
				FOR SELECT `clubs`.`id`
					FROM `clubs`
					WHERE FIND_IN_SET(`clubs`.`display_id`, `clubs`);

			DECLARE `old` CURSOR FOR SELECT * FROM `temporary`;

			DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
				BEGIN
					DROP TABLE `temporary`;
					RESIGNAL;
				END;

			INSERT INTO `temporary` (`club`)
				SELECT `club_member`.`club`
					FROM `club_member`
					WHERE `club_member`.`member`=`member`;

			OPEN `new`;
			BEGIN
				DECLARE `count` tinyint unsigned DEFAULT 0;
				DECLARE `club` tinyint unsigned;
				DECLARE EXIT HANDLER FOR NOT FOUND
					IF `count` < `clubs_number` THEN
						SIGNAL SQLSTATE '45000';
					END IF;

				LOOP
					FETCH `new` INTO `club`;
					DELETE FROM `temporary` WHERE `temporary`.`club`=`club`;
					IF ROW_COUNT() <= 0 THEN
						INSERT INTO `club_member` (`club`, `member`)
							VALUES (`club`, `member`);
					END IF;
					SET `count` = `count` + 1;
				END LOOP;
			END;
			CLOSE `new`;

			OPEN `old`;
			BEGIN
				DECLARE `club` tinyint unsigned;
				DECLARE EXIT HANDLER FOR NOT FOUND BEGIN END;

				LOOP
					FETCH `old` INTO `club`;
					DELETE
						FROM `club_member`
						WHERE
							`club_member`.`club`=`club` AND
							`club_member`.`member`=`member`;
				END LOOP;
			END;
			CLOSE `old`;
		END;

		DROP TABLE `temporary`;
		COMMIT;
	END$

CREATE OR REPLACE PROCEDURE `update_officer` (
	`operator` varchar(255) CHARACTER SET ascii,
	`display_id` varchar(255) CHARACTER SET ascii,
	`name` varchar(63) CHARACTER SET utf8mb4,
	`member` varchar(255) CHARACTER SET ascii,
	`scope` set('management', 'privacy') CHARACTER SET ascii)
	MODIFIES SQL DATA
	COMMENT 'TODO'
	BEGIN
		IF NOT FIND_IN_SET('management', `scope`) AND
			`officer_depending_for_management`(`operator`, `display_id`)
		THEN
			SIGNAL SQLSTATE '45000';
		END IF;

		UPDATE `officers`
			SET
				`name`=IFNULL(`name`, `officers`.`name`),
				`member`=IFNULL(`member`, `officers`.`member`),
				`scope`=IFNULL(`scope`, `officers`.`scope`)
			WHERE `officers`.`display_id`=`display_id`;
	END$

CREATE OR REPLACE PROCEDURE `update_party` (
	`name` varchar(63) CHARACTER SET utf8mb4,
	`creator` varchar(255) CHARACTER SET ascii,
	`start` datetime,
	`end` datetime,
	`place` varchar(63) CHARACTER SET utf8mb4,
	`due` datetime,
	`inviteds` varchar(63) CHARACTER SET utf8mb4,
	`invited_ids` varchar(65535) CHARACTER SET ascii,
	`inviteds_number` smallint unsigned,
	`details` varchar(63) CHARACTER SET utf8mb4)
	MODIFIES SQL DATA
	COMMENT 'TODO'
	BEGIN
		DECLARE `party` smallint unsigned;

		DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
			BEGIN
				ROLLBACK;
				RESIGNAL;
			END;

		SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
		START TRANSACTION;

		DO LAST_INSERT_ID(0);

		UPDATE `parties`
			SET
				`id`=LAST_INSERT_ID(`parties`.`id`),
				`start`=IFNULL(`start`, `members`.`start`),
				`end`=IFNULL(`end`, `members`.`end`),
				`place`=IFNULL(`place`, `members`.`place`),
				`due`=IFNULL(`due`, `members`.`due`),
				`inviteds`=IFNULL(`inviteds`, `members`.`inviteds`),
				`details`=IFNULL(`details`, `members`.`details`)
			WHERE `members`.`display_id`=`display_id`;

		SET `party`=LAST_INSERT_ID();

		IF `party`=0 THEN
			SIGNAL SQLSTATE '45000';
		END IF;

		CREATE TEMPORARY TABLE `temporary` (`member` tinyint unsigned NOT NULL PRIMARY KEY);
		BEGIN
			DECLARE `new` CURSOR
				FOR SELECT `members`.`id`
					FROM `members`
					WHERE FIND_IN_SET(`members`.`display_id`, `invited_ids`);

			DECLARE `old` CURSOR FOR SELECT * FROM `temporary`;

			DECLARE EXIT HANDLER FOR SQLEXCEPTION, NOT FOUND
				BEGIN
					DROP TABLE `temporary`;
					RESIGNAL;
				END;

			INSERT INTO `temporary` (`member`)
				SELECT `attendances`.`member`
					FROM `attendances`
					WHERE `attendances`.`party`=`party`;

			OPEN `new`;
			BEGIN
				DECLARE `count` tinyint unsigned DEFAULT 0;
				DECLARE `member` tinyint unsigned;
				DECLARE EXIT HANDLER FOR NOT FOUND
					IF `count` < `inviteds_number` THEN
						SIGNAL SQLSTATE '45000';
					END IF;

				LOOP
					FETCH `new` INTO `member`;
					DELETE FROM `temporary` WHERE `temporary`.`member`=`member`;
					IF ROW_COUNT() <= 0 THEN
						INSERT INTO `attendances` (`party`, `member`)
							VALUES (`party`, `member`);
					END IF;
					SET `count` = `count` + 1;
				END LOOP;
			END;
			CLOSE `new`;

			OPEN `old`;
			BEGIN
				DECLARE `member` tinyint unsigned;
				DECLARE EXIT HANDLER FOR NOT FOUND BEGIN END;

				LOOP
					FETCH `old` INTO `member`;
					DELETE
						FROM `attendances`
						WHERE
							`attendances`.`party`=`party` AND
							`attendances`.`member`=`member`;
				END LOOP;
			END;
			CLOSE `old`;
		END;

		DROP TABLE `temporary`;
		COMMIT;
	END$

DELIMITER ;

SET sql_mode=@saved_sql_mode;
SET foreign_key_checks=@saved_foreign_key_checks;
SET unique_checks=@saved_unique_checks;
SET time_zone=@saved_time_zone;
SET collation_connection=@saved_collation_connection;
SET character_set_results=@saved_character_set_results;
SET character_set_client=@saved_character_set_client;

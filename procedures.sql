DELIMITER !
CREATE PROCEDURE brief_club(IN in_display_id varchar(255) CHARACTER SET ascii)
	BEGIN
		DECLARE selected_chief smallint unsigned;
		DECLARE selected_id tinyint unsigned;
		DECLARE selected_name varchar(64);
		DECLARE EXIT HANDLER FOR NOT FOUND BEGIN END;
		SELECT chief, id, name FROM clubs WHERE display_id = in_display_id
			INTO selected_chief, selected_id, selected_name;

		SELECT selected_name;
		SELECT display_id, mail, nickname, realname, tel FROM members WHERE id = selected_chief;

		member_scope: BEGIN
			DECLARE done bool DEFAULT FALSE;
			DECLARE member_id tinyint unsigned;
			DECLARE member_cursor CURSOR FOR SELECT member FROM club_member WHERE club = selected_id;
			DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

			OPEN member_cursor;

			member_loop: LOOP
				FETCH FROM member_cursor INTO member_id;

				IF done THEN
					LEAVE member_loop;
				END IF;

				SELECT entrance,display_id, nickname, realname FROM members WHERE id = member_id;
			END LOOP member_loop;

			CLOSE member_cursor;
		END member_scope;
	END!

CREATE PROCEDURE insert_member(
	IN in_affiliation varchar(63) CHARACTER SET utf8mb4,
	IN in_display_id varchar(255) CHARACTER SET ascii,
	IN in_entrance year,
	IN in_gender varchar(63) CHARACTER SET utf8mb4,
	IN in_mail varchar(255) CHARACTER SET ascii,
	IN in_nickname varchar(63) CHARACTER SET utf8mb4,
	IN in_realname varchar(63) CHARACTER SET utf8mb4,
	IN in_tel varchar(63) CHARACTER SET utf8mb4)
	BEGIN
		INSERT INTO members (affiliation, display_id, entrance, gender, mail, nickname, realname, tel)
			VALUES (in_affiliation, in_display_id, in_entrance, in_gender, in_mail, in_nickname, in_realname, in_tel);
		SELECT LAST_INSERT_ID();
	END!

DELIMITER ;

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

package db

/*
	Check whether they are used with the following command (replace LIST):
	for q in LIST; do
		if [ "`git grep $q -- backend/db ':!backend/db/stmts.go'`" = '' ]; then
			echo $q;
		fi;
	done;
*/
const (
	stmtCallDeleteOfficer = iota
	stmtCallInsertMail
	stmtCallInsertParty
	stmtCallUpdateMail
	stmtCallUpdateMember
	stmtCallUpdateOfficer
	stmtCallUpdateParty
	stmtConfirmMember
	stmtCountMembers
	stmtDeclareMemberOB
	stmtDeleteClub
	stmtDeleteMail
	stmtDeleteMember
	stmtDeleteParty
	stmtInsertClub
	stmtInsertMember
	stmtInsertOfficer
	stmtSelectAttendancesByInternalParty
	stmtSelectAttendancesByMember
	stmtSelectClubByID
	stmtSelectClubIDInternalMembers
	stmtSelectClubInternalIDMemberID
	stmtSelectClubNameByID
	stmtSelectClubs
	stmtSelectClubsByInternalMember
	stmtSelectMailBySubject
	stmtSelectMails
	stmtSelectMemberByID
	stmtSelectMemberGraphByID
	stmtSelectMemberIDMails
	stmtSelectMemberIDsByInternalClub
	stmtSelectMemberInternalIDPasswordByID
	stmtSelectMemberInternalIDNicknameByID
	stmtSelectMemberNicknameByID
	stmtSelectMemberPasswordByID
	stmtSelectMemberRoles
	stmtSelectMembers
	stmtSelectOfficerByID
	stmtSelectOfficerIDByMemberID
	stmtSelectOfficerIDNames
	stmtSelectOfficerNameByID
	stmtSelectOfficerScopeByInternalMember
	stmtSelectOfficers
	stmtSelectParties
	stmtSelectParty
	stmtSelectRecipientsByInternalMail
	stmtUpdateAttendance
	stmtUpdateClub
	stmtUpdateMemberPassword

	stmtNumber
)

var stmtQueries = [...]string{
	stmtCallDeleteOfficer:                  "CALL `delete_officer`(?, ?)",
	stmtCallInsertMail:                     "CALL `insert_mail`(?, ?, ?, ?, ?, ?)",
	stmtCallInsertParty:                    "CALL `insert_party`(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	stmtCallUpdateMail:                     "CALL `update_mail(?, ?, ?, ?, ?, ?, ?)",
	stmtCallUpdateMember:                   "CALL `update_member`(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	stmtCallUpdateOfficer:                  "CALL `update_officer`(?, ?, ?, ?, ?)",
	stmtCallUpdateParty:                    "CALL `update_party`(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	stmtConfirmMember:                      "UPDATE `members` SET `flags`=`flags`|1 WHERE `display_id`=?",
	stmtCountMembers:                       "SELECT COUNT(*) FROM `members` WHERE `nickname` LIKE ? AND `realname` LIKE ? AND `entrance`=? IS NOT FALSE AND FIND_IN_SET('ob', `flags`)=? IS NOT FALSE",
	stmtDeclareMemberOB:                    "UPDATE `members` SET `flags`=`flags`|2 WHERE `display_id`=?",
	stmtDeleteClub:                         "DELETE FROM `clubs` WHERE `display_id`=?",
	stmtDeleteMail:                         "DELETE FROM `mails` WHERE `subject`=?",
	stmtDeleteMember:                       "DELETE FROM `members` WHERE `display_id`=?",
	stmtDeleteParty:                        "DELETE `parties` FROM `parties` JOIN `members` ON `parties`.`creator`=`members`.`display_id` WHERE `parties`.`name`=? AND `members`.`display_id`=?",
	stmtInsertClub:                         "INSERT `clubs` (`display_id`, `name`, `chief`) SELECT ?, ?, `id` FROM `members` WHERE `display_id`=?",
	stmtInsertMember:                       "INSERT `members` (`display_id`, `mail`, `nickname`) VALUES (?, ?, ?)",
	stmtInsertOfficer:                      "INSERT `officers` (`display_id`, `name`, `scope`, `member`) SELECT ?, ?, ?, `id` FROM `members` WHERE `display_id`=?",
	stmtSelectAttendancesByInternalParty:   "SELECT `members`.`display_id`, CAST(`attendances`.`attendance` as int) FROM `members` JOIN `attendances` ON `members`.`id`=`attendances`.`member` WHERE `attendances`.`party`=?",
	stmtSelectAttendancesByMember:          "SELECT `attendances`.`party`, CAST(`attendances`.`attendance` as int) FROM `members` JOIN `attendances` ON `members`.`id`=`attendances`.`member` WHERE `members`.`display_id`=?",
	stmtSelectClubByID:                     "SELECT `clubs`.`id`, `clubs`.`name`, `members`.`display_id` FROM `clubs` JOIN `members` ON `clubs`.`chief`=`members`.`id` WHERE `clubs`.`display_id`=?",
	stmtSelectClubIDInternalMembers:        "SELECT `clubs`.`display_id`, `club_member`.`member` FROM `clubs` JOIN `club_member` ON `clubs`.`id`=`club_member`.`club`",
	stmtSelectClubInternalIDMemberID:       "SELECT `club_member`.`club`, `members`.`display_id` FROM `club_member` JOIN `members` ON `club_member`.`member`=`members`.`id`",
	stmtSelectClubNameByID:                 "SELECT `name` FROM `clubs` WHERE `display_id`=?",
	stmtSelectClubs:                        "SELECT `clubs`.`id`, `clubs`.`display_id`, `clubs`.`name`, `members`.`display_id` FROM `clubs` JOIN `members` ON `clubs`.`chief`=`members`.`id`",
	stmtSelectClubsByInternalMember:        "SELECT `clubs`.`chief`, `clubs`.`display_id` FROM `club_member` JOIN `clubs` ON `club_member`.`club`=`clubs`.`id` WHERE `club_member`.`member`=?",
	stmtSelectMailBySubject:                "SELECT `mails`.`id`, `mails`.`date`, `members`.`display_id`, `mails`.`to`, `mails`.`body` FROM `mails` LEFT JOIN `members` ON `mails`.`from`=`members`.`id` WHERE `mails`.`subject`=?",
	stmtSelectMails:                        "SELECT `mails`.`date`, `members`.`display_id`, `mails`.`to`, `mails`.`subject` FROM `mails` LEFT JOIN `members` ON `mails`.`from`=`members`.`id`",
	stmtSelectMemberByID:                   "SELECT `id`, `affiliation`, `entrance`, CAST(`flags` as int), `gender`, `mail`, `nickname`, `realname`, `tel` FROM `members` WHERE `display_id`=?",
	stmtSelectMemberGraphByID:              "SELECT `gender`, `nickname` FROM `members` WHERE `display_id`=?",
	stmtSelectMemberIDsByInternalClub:        "SELECT `members`.`display_id` FROM `club_member` JOIN `members` ON `club_member`.`member`=`members`.`id` WHERE `club`=?",
	stmtSelectMemberIDMails:                "SELECT `display_id`, `mail` FROM `members`",
	stmtSelectMemberInternalIDPasswordByID: "SELECT `id`, `password` FROM `members` WHERE `display_id`=?",
	stmtSelectMemberInternalIDNicknameByID: "SELECT `id`, `nickname` FROM `members` WHERE `display_id`=?",
	stmtSelectMemberNicknameByID:           "SELECT `nickname` FROM `members` WHERE `display_id`=?",
	stmtSelectMemberPasswordByID:           "SELECT `password` FROM `members` WHERE `display_id`=?",
	stmtSelectMemberRoles:                  "SELECT `display_id`, CAST(`flags` as int), `id`, `nickname` FROM `members`",
	stmtSelectMembers:                      "SELECT `affiliation`, `display_id`, `entrance`, CAST(`flags` as int), `nickname`, `realname` FROM `members`",
	stmtSelectOfficerByID:                  "SELECT `officers`.`name`, `officers`.`scope`, `members`.`display_id` FROM `officers` JOIN `members` ON `officers`.`member`=`members`.`id` WHERE `officers`.`display_id`=?",
	stmtSelectOfficerIDByMemberID:          "SELECT `display_id` FROM `officers` WHERE `member`=?",
	stmtSelectOfficerIDNames:               "SELECT `display_id`, `name` FROM `officers`",
	stmtSelectOfficerNameByID:              "SELECT `name` FROM `officers` WHERE `display_id`=?",
	stmtSelectOfficerScopeByInternalMember: "SELECT `scope` FROM `officers` WHERE `member`=?",
	stmtSelectOfficers:                     "SELECT `officers`.`display_id`, `officers`.`name`, `members`.`display_id` FROM `officers` JOIN `members` ON `officers`.`member`=`members`.`id`",
	stmtSelectParties:                      "SELECT `parties`.`id`, `parties`.`name`, `members`.`display_id`, `parties`.`start`, `parties`.`end`, `parties`.`place`, `parties`.`inviteds`, `parties`.`due` FROM `parties` LEFT JOIN `members` ON `parties`.`creator`=`members`.`id`",
	stmtSelectParty:                        "SELECT `parties`.`id`, `members`.`display_id`, `parties`.`start`, `parties`.`end`, `parties`.`place`, `parties`.`inviteds`, `parties`.`due`, `parties`.`details` FROM `parties` LEFT JOIN `members` ON `parties`.`creator`=`members`.`id` WHERE `parties`.`name`=?",
	stmtSelectRecipientsByInternalMail:     "SELECT `members`.`display_id` FROM `members` JOIN `recipients` ON `members`.`id`=`recipients`.`member` WHERE `recipients`.`mail`=?",
	stmtUpdateAttendance:                   "UPDATE `attendances` JOIN `parties` ON `attendances`.`party`=`parties`.`id` JOIN `members` ON `attendances`.`member`=`members`.`id` SET `attendances`.`attendance`=? WHERE `parties`.`name`=? AND `members`.`display_id`=?",
	stmtUpdateClub:                         "UPDATE `clubs` SET `name`=IFNULL(@name, `name`), `chief`=IF(@chief, (SELECT `id` FROM `members` WHERE `display_id`=@chief), `chief`) WHERE `display_id`=@id",
	stmtUpdateMemberPassword:               "UPDATE `members` SET `password`=? WHERE `display_id`=?",
}

func (db *DB) prepareStmts() error {
	for index, query := range stmtQueries {
		var err error

		db.stmts[index], err = db.sql.Prepare(query)
		if err != nil {
			return err
		}
	}

	return nil
}

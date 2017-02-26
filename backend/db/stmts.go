package db

const (
	stmtConfirmMember = iota
	stmtCountMembers
	stmtDeclareMemberOB
	stmtDeleteClubMemberByInternal
	stmtDeleteMember
	stmtInsertInternalAttendance
	stmtInsertInternalClubMember
	stmtInsertMember
	stmtInsertParty
	stmtSelectAttendancesByMember
	stmtSelectClub
	stmtSelectClubIDsByInternalMember
	stmtSelectClubInternalByInternalMember
	stmtSelectClubName
	stmtSelectClubNames
	stmtSelectClubs
	stmtSelectClubsByInternalMember
	stmtSelectInternalClubs
	stmtSelectMember
	stmtSelectMemberGraph
	stmtSelectMemberID
	stmtSelectMemberIDPassword
	stmtSelectMemberInternalIDMails
	stmtSelectMemberMails
	stmtSelectMemberNickname
	stmtSelectMemberOfficer
	stmtSelectMemberPassword
	stmtSelectMemberRoles
	stmtSelectMembers
	stmtSelectMembersByClub
	stmtSelectOfficer
	stmtSelectOfficerName
	stmtSelectOfficerScopeByInternalMember
	stmtSelectOfficers
	stmtSelectParties
	stmtSelectPartyNames
	stmtUpdateAttendance
	stmtUpdateMemberPassword

	stmtNumber
)

var stmtQueries = [...]string{
	stmtConfirmMember:                     `UPDATE members SET flags=flags|1 WHERE display_id=?`,
	stmtCountMembers:                      `SELECT COUNT(id) FROM members WHERE nickname LIKE ? AND realname LIKE ? AND entrance=? IS NOT FALSE AND FIND_IN_SET('ob', flags) IS NOT FALSE`,
	stmtDeclareMemberOB:                   `UPDATE members SET flags=flags|1 WHERE display_id=?`,
	stmtDeleteClubMemberByInternal:        `DELETE FROM club_member WHERE FIND_IN_SET(id, ?)`,
	stmtDeleteMember:                      `DELETE FROM members WHERE display_id=?`,
	stmtInsertInternalAttendance:          `INSERT attendances (party, member) VALUES (?, ?)`,
	stmtInsertInternalClubMember:          `INSERT club_member (club, member) VALUES (?, ?)`,
	stmtInsertMember:                      `INSERT members (display_id, mail, nickname) VALUES (?, ?, ?)`,
	stmtInsertParty:                       `INSERT parties (name, start, end, place, due, inviteds, details) VALUES (?, ?, ?, ?, ?, ?, ?)`,
	stmtSelectAttendancesByMember:         `SELECT attendances.party, attendances.attending FROM members JOIN attendances ON members.id=attendances.member WHERE members.display_id=?`,
	stmtSelectClub:                        `SELECT clubs.id, clubs.name, members.display_id, members.mail, members.nickname, members.realname, members.tel FROM clubs JOIN members ON clubs.chief=members.id WHERE clubs.display_id=?`,
	stmtSelectClubIDsByInternalMember:      `SELECT clubs.display_id FROM club_member JOIN clubs ON club_member.club=clubs.id WHERE club_member.member=?`,
	stmtSelectClubInternalByInternalMember: `SELECT id, club FROM club_member WHERE member=?`,
	stmtSelectClubName:                     `SELECT name FROM clubs WHERE display_id=?`,
	stmtSelectClubNames:                    `SELECT display_id, name FROM clubs`,
	stmtSelectClubs:                        `SELECT clubs.display_id, clubs.name, members.display_id, members.mail, members.nickname, members.realname, members.tel FROM clubs JOIN members ON clubs.chief=members.id`,
	stmtSelectClubsByInternalMember:        `SELECT clubs.chief, clubs.display_id, clubs.name FROM club_member JOIN clubs ON club_member.club=clubs.id WHERE club_member.member=?`,
	stmtSelectInternalClubs:                `SELECT id FROM clubs WHERE FIND_IN_SET(display_id, ?)`,
	stmtSelectMember:                       `SELECT id, affiliation, entrance, flags, gender, mail, nickname, realname, tel FROM members WHERE display_id=?`,
	stmtSelectMemberGraph:                  `SELECT gender, nickname FROM members WHERE display_id=?`,
	stmtSelectMemberID:                     `SELECT id FROM members WHERE display_id=?`,
	stmtSelectMemberIDPassword:             `SELECT id, password FROM members WHERE display_id=?`,
	stmtSelectMemberInternalIDMails:        `SELECT id,mail FROM members WHERE FIND_IN_SET(display_id, ?)`,
	stmtSelectMemberMails:                  `SELECT mail FROM members WHERE FIND_IN_SET(display_id, ?)`,
	stmtSelectMemberNickname:               `SELECT nickname FROM members WHERE display_id=?`,
	stmtSelectMemberOfficer:                `SELECT display_id, name FROM officers WHERE member=?`,
	stmtSelectMemberPassword:               `SELECT password FROM members WHERE display_id=?`,
	stmtSelectMemberRoles:                  `SELECT display_id,flags,id,nickname FROM members`,
	stmtSelectMembers:                      `SELECT affiliation, display_id, entrance, flags, nickname, realname FROM members`,
	stmtSelectMembersByClub:                `SELECT members.entrance, members.display_id, members.nickname, members.realname FROM club_member JOIN members ON club_member.member=members.id WHERE club=?`,
	stmtSelectOfficer:                      `SELECT officers.name, officers.scope, members.display_id, members.mail, members.nickname, members.realname, members.tel FROM officers JOIN members ON officers.member=members.id WHERE officers.display_id=?`,
	stmtSelectOfficerName:                  `SELECT name FROM officers WHERE display_id=?`,
	stmtSelectOfficerScopeByInternalMember: `SELECT scope FROM officers WHERE member=?`,
	stmtSelectOfficers:                     `SELECT officers.display_id, officers.name, members.display_id, members.mail, members.nickname, members.realname, members.tel FROM officers JOIN members ON officers.member=members.id`,
	stmtSelectParties:                      `SELECT id, name, start, end, place, inviteds, due FROM parties`,
	stmtSelectPartyNames:                   `SELECT name FROM parties`,
	stmtUpdateAttendance:                   `UPDATE attendances JOIN parties ON attendances.party=parties.id JOIN members ON attendances.member=members.id SET attendances.attending=? WHERE parties.name=? AND members.display_id=?`,
	stmtUpdateMemberPassword:               `UPDATE members SET password=? WHERE display_id=?`,
}

func (db *DB) prepareStmts() error {
	for index, query := range stmtQueries {
		var prepareError error
		db.stmts[index], prepareError = db.sql.Prepare(query)
		if prepareError != nil {
			return prepareError
		}
	}

	return nil
}

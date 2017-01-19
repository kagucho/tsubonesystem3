package db

const (
	stmtCountMembers = iota
	stmtDeclareOB
	stmtDeleteMember
	stmtInsertClubMember
	stmtInsertClubMemberInternal
	stmtInsertMember
	stmtSelectBriefMemberInternal
	stmtSelectClub
	stmtSelectClubInternal
	stmtSelectClubMembers
	stmtSelectClubName
	stmtSelectClubNames
	stmtSelectClubs
	stmtSelectMember
	stmtSelectMemberClubIDs
	stmtSelectMemberClubs
	stmtSelectMemberGraph
	stmtSelectMemberID
	stmtSelectMemberIDPassword
	stmtSelectMemberOfficer
	stmtSelectMemberPassword
	stmtSelectMembers
	stmtSelectOfficer
	stmtSelectOfficerMemberInternal
	stmtSelectOfficerName
	stmtSelectOfficerScopeInternal
	stmtSelectOfficers

	stmtNumber
)

var stmtQueries = [...]string{
	stmtCountMembers:                `SELECT COUNT(id) FROM members WHERE nickname LIKE ? AND realname LIKE ? AND entrance=? IS NOT FALSE AND ob=? IS NOT FALSE`,
	stmtDeclareOB:                   `UPDATE members SET ob=1 WHERE display_id=?`,
	stmtDeleteMember:                `DELETE FROM members WHERE id=?`,
	stmtInsertClubMember:            `INSERT club_member(club,member)VALUES((SELECT id FROM clubs WHERE display_id=?),?)`,
	stmtInsertClubMemberInternal:    `INSERT club_member(club,member)VALUES(?,?)`,
	stmtInsertMember:                `INSERT members(display_id,mail,nickname)VALUES(?,?,?)`,
	stmtSelectBriefMemberInternal:   `SELECT entrance,display_id,nickname,realname FROM members WHERE id=?`,
	stmtSelectClub:                  `SELECT chief,id,name FROM clubs WHERE display_id=?`,
	stmtSelectClubInternal:          `SELECT chief,display_id,name FROM clubs WHERE id=?`,
	stmtSelectClubMembers:           `SELECT member FROM club_member WHERE club=?`,
	stmtSelectClubName:              `SELECT name FROM clubs WHERE display_id=?`,
	stmtSelectClubNames:             `SELECT display_id,name FROM clubs`,
	stmtSelectClubs:                 `SELECT chief,display_id,name FROM clubs`,
	stmtSelectMember:                `SELECT id,affiliation,entrance,gender,mail,nickname,ob,realname,tel FROM members WHERE display_id=?`,
	stmtSelectMemberClubIDs:         `SELECT id,club FROM club_member WHERE member=?`,
	stmtSelectMemberClubs:           `SELECT club FROM club_member WHERE member=?`,
	stmtSelectMemberGraph:           `SELECT gender,nickname FROM members WHERE display_id=?`,
	stmtSelectMemberID:              `SELECT id FROM members WHERE display_id=?`,
	stmtSelectMemberIDPassword:      `SELECT id,password FROM members WHERE display_id=?`,
	stmtSelectMemberOfficer:         `SELECT display_id,name FROM officers WHERE member=?`,
	stmtSelectMemberPassword:        `SELECT password FROM members WHERE display_id=?`,
	stmtSelectMembers:               `SELECT affiliation,entrance,display_id,nickname,ob,realname FROM members`,
	stmtSelectOfficer:               `SELECT member,name,scope FROM officers WHERE display_id=?`,
	stmtSelectOfficerMemberInternal: `SELECT display_id,mail,nickname,realname,tel FROM members WHERE id=?`,
	stmtSelectOfficerName:           `SELECT name FROM officers WHERE display_id=?`,
	stmtSelectOfficerScopeInternal:  `SELECT scope FROM officers WHERE member=?`,
	stmtSelectOfficers:              `SELECT display_id,member,name FROM officers`,
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

/**
	@file table.js implements tables.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/table */

import * as modal from "../modal";
import * as url from "../url";
import large from "../large";
import memberModal from "./member/modal";

/**
	showMember shows a member.
	@private
	@param {?String} id - The ID. If it is null, it will be the form to
	create a new member. Otherwise, it will describe the member identified
	with the ID.
*/
function showMember(id) {
	if (large()) {
		this.modal = modal.add(memberModal(id));

		return false;
	}

	return true;
}

/**
	module:private/component/table is a class which contains the common
	implementation of tables.
	@implements external:Mithril~Component
*/
export default class Table {
	/**
		constructor returns a new module:private/component/table#Table.
		@param {!module:private/component/table~TableView} table
		- A function which returns the view of the table.
		@returns {!module:private/component/table#Table} A new
		instance.
	*/
	constructor(table) {
		this.table = table;
		Object.freeze(this);
	}

	view(node) {
		return m("div", {className: "table-responsive"},
			m("table", {className: "table"},
				this.table(node.attrs.members, member => m("a", {
					href:    "#!member?id=" + member.id,
					onclick: showMember.bind(this,
						member.id),
				}, member.nickname))
			)
		);
	}
}

/**
	A function which returns the view of a table.
	@callback module:private/component/table~TableView
	@param {!*} entries - The entries to draw.
	@param {!module:private/component/table~NicknameView} - A function
	which returns the view of a member.
	@returns {!external:Mithril~Children} The view of a table.
*/

/**
	A function which returns the view of a nickname.
	@callback module:private/component/table~NicknameView
	@param {!*} entry - The entry whose nickname is to
	be drawn.
	@returns {!external:Mithril~Children} The view of the nickname.
*/

/**
	members is a component to draw tables of members.
	@type external:Mithril~Component
*/
export const members = new Table((entries, nicknameView) => [
	m("thead", m("tr", {style: {backgroundColor: "#d9edf7"}},
		m("th", "ニックネーム"),
		m("th", "名前"),
		m("th", "入学年度")
	)), m("tbody", entries.map(entry => m("tr", {key: entry.id},
		m("td", nicknameView(entry)),
		m("td", entry.realname),
		m("td", entry.entrance)
	))),
]);

/**
	officers is a component to draw tables of officers.
	@type external:Mithril~Component
*/
export const officers = new Table((entries, nicknameView) => [
	m("thead", m("tr",
		m("th", "ニックネーム"),
		m("th", "名前"),
		m("th", "メールアドレス"),
		m("th", "電話番号")
	)), m("tbody", entries.map(entry => m("tr", {key: entry.id},
		m("td", nicknameView(entry)),
		m("td", entry.realname),
		m("td", m("a", {href: url.mailto(entry.mail)}, entry.mail)),
		m("td", m("a", {href: url.tel(entry.tel)}, entry.tel))
	))),
]);

/**
	State is the state of the table.
	@typedef {module:private/component/table}
		module:private/component/table~State
	@property {?module:private/modal~Node} modal - The node of the modal
	dialog of a member in the list of modal dialogs.
*/

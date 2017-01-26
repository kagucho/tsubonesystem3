/**
	@file table.js implements tables.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/table */

import * as memberModal from "./member/modal";
import * as url from "../url";
import large from "../large";

/**
	Common is a class which contains the common implementation of tables.
	@implements external:Mithril~Component
*/
class Common {
	/**
		constructor returns a new instance of Common class.
		@param {!module:private/components/table~tableView} table
		- A function which returns the view of the table.
		@returns {!module:private/components/table/Common} A new
		instance.
	*/
	constructor(table) {
		this.table = table;
		Object.freeze(this);
	}

	controller() {
		return {
			showMember(id) {
				if (large()) {
					this.rendering = id;

					return false;
				}
			},

			hideMember() {
				delete this.rendering;
			},
		};
	}

	view(control, attributes) {
		return m("div", {className: "table-responsive"},
			m("table", {className: "table"},
				this.table(attributes.members, member => m("a", {
					href:    "#!member?id=" + member.id,
					onclick: control.showMember.bind(control, member.id),
				}, member.nickname))
			), control.rendering && m(memberModal, {
				id:          control.rendering,
				onhidden:    control.hideMember.bind(control),
				onloadstart: attributes.onloadstart,
			}));
	}
}

/**
	A function which returns the view of a table.
	@callback module:private/components/table~tableView
	@param {!external:ES.Object} entries - The entries to draw.
	@param {!module:private/components/table~nicknameView} - A function
	which returns the view of a member.
	@returns {!external:Mithril~Children} The view of a table.
*/

/**
	A function which returns the view of a nickname.
	@callback module:private/components/table~nicknameView
	@param {!external:ES.Object} entry - The entry whose nickname is to
	be drawn.
	@returns {!external:Mithril~Children} The view of the nickname.
*/

/**
	members is a component to draw tables of members.
	@type external:Mithril~Component
*/
export const members = new Common((entries, nicknameView) => [
	m("thead", m("tr", {style: {backgroundColor: "#d9edf7"}},
		m("th", "ニックネーム"),
		m("th", "名前"),
		m("th", "入学年度")
	)), m("tbody", entries.map(entry => m("tr",
		m("td", nicknameView(entry)),
		m("td", entry.realname),
		m("td", entry.entrance)
	))),
]);

/**
	officers is a component to draw tables of officers.
	@type external:Mithril~Component
*/
export const officers = new Common((entries, nicknameView) => [
	m("thead", m("tr",
		m("th", "ニックネーム"),
		m("th", "名前"),
		m("th", "メールアドレス"),
		m("th", "電話番号")
	)), m("tbody", entries.map(entry => m("tr",
		m("td", nicknameView(entry)),
		m("td", entry.realname),
		m("td", m("a", {href: url.mailto(entry.mail)}, entry.mail)),
		m("td", m("a", {href: url.tel(entry.tel)}, entry.tel))
	))),
]);

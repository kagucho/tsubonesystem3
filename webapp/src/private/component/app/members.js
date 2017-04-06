/**
	@file members.js implements members component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/members */

/**
	module:private/component/app/members is a component to show members.
	@name module:private/component/app/members
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as modal from "../../modal";
import * as progress from "../../progress";
import Table from "../table";
import client from "../../client";
import large from "../../large";
import memberModal from "../member/modal";

/**
	addMember prompts the user to add a new member.
	@private
	@this module:private/component/app/members
	@returns {!Boolean} false if the user agent should not leave this page.
*/
function addMember() {
	if (large()) {
		this.memberModal = modal.add(memberModal());

		return false;
	}

	return true;
}

/**
	setTableState sets the state of the table.
	@private
	@this module:private/component/app/members
	@param {!external:Mithril~Node} node - The node of the table.
	@returns {Undefined}
*/
function setTableState(node) {
	this.tableState = node.state;
}

/**
	filter filters the member to show.
	@private
	@this module:private/component/app/members
	@returns {Undefined}
*/
function filter() {
	this.count = 0;

	for (const member of this.members) {
		if ((!this.param.entrance ||
				member.entrance == this.param.entrance) &&
			(!this.param.nickname ||
				member.nickname.includes(this.param.nickname)) &&
			(!this.param.ob ||
				this.param.ob == "1" && member.ob ||
				this.param.ob == "0" && !member.ob) &&
			(!this.param.realname ||
				member.realname.includes(this.param.realname))) {
			member.display = true;
			this.count++;
		} else {
			member.display = false;
		}
	}
}

/**
	update updates the value identified by the given key.
	@private
	@this module:private/component/app/members
	@param {!String} key - The key to identify the value.
	@param {!String} value - The new value to be set.
	@returns {Undefined}
*/
function update(key, value) {
	if (value) {
		this.param[key] = value;
	} else {
		delete this.param[key];
	}

	let url = "#!members";
	if (!$.isEmptyObject(this.param)) {
		url += "?" + m.buildQueryString(this.param);
	}

	history.pushState(null, null, url);
	filter.call(this);
}

/**
	tableView returns the virtual DOM nodes of a table.
	@private
	@param {!Array} entries - Entries of the table.
	@param {!module:private/component/table~NicknameView} nicknameView
	- The view of nicknames.
	@returns {external:Mithril~Children}
*/
const tableView = (entries, nicknameView) => [
	m("thead", m("tr", {style: {backgroundColor: "#d9edf7"}},
		m("th", "ニックネーム"),
		m("th", "名前"),
		m("th", "入学年度")
	)), m("tbody", entries.map(entry => m("tr", {
		"aria-hidden": (!entry.display).toString(),
		key:           entry.id,
		style:         {display: entry.display ? "table-row" : "none"},
	},
		m("td", nicknameView(entry)),
		m("td", entry.realname),
		m("td", entry.entrance)
	))),
];

export function oninit() {
	this.param = m.route.param();

	this.stream = client.mapMembers(promise => {
		const loadingProgress = progress.add({
			"aria-describedby": "component-app-members-loading",
			value:              0,
		});

		promise.then(members => {
			this.members = Array.from(members());
			filter.call(this);
			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.error = client.error(error);
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	});

	this.table = new Table(tableView);
}

export function onbeforeremove() {
	if (this.memberModal) {
		this.memberModal.remove();
	}

	if (this.stream) {
		this.stream.end(true);
	}

	if (this.tableState && this.tableState.modal) {
		this.tableState.modal.remove();
	}
}

export function view() {
	const promise = this.stream();

	return [
		m(container,
			m("div", {className: "container"},
				m("h1", "Members"),
				m("div", {
					role:  "search",
					style: {
						backgroundColor: "#f0f8ff",
						borderRadius:    "1rem",
						margin:          "1rem",
						padding:         "1rem",
					},
				},
					m("div", {style: {display: "table"}}, [
						{
							label: m("label", {
								className: "control-label",
								htmlFor:   "component-app-members-nickname",
							}, "ニックネーム"),
							input: m("input", {
								className: "form-control",
								id:        "component-app-members-nickname",
								maxlength: "63",
								oninput:   m.withAttr("value",
									update.bind(this, "nickname")),
								placeholder: "Nickname",
								type:        "search",
								value:       this.param.nickname || "",
							}),
						}, {
							label: m("label", {
								className: "control-label",
								htmlFor:   "component-app-members-realname",
							}, "名前"),
							input: m("input", {
								className: "form-control",
								id:        "component-app-members-realname",
								maxlength: "63",
								oninput:   m.withAttr("value",
									update.bind(this, "realname")),
								placeholder: "Name",
								type:        "search",
								value:       this.param.realname || "",
							}),
						}, {
							label: m("label", {
								className: "control-label",
								htmlFor:   "component-app-members-entrance",
							}, "入学年度"),
							input: m("input", {
								className: "form-control",
								id:        "component-app-members-entrance",
								max:       "2155",
								min:       "1901",
								oninput:   m.withAttr("value",
									update.bind(this, "entrance")),
								placeholder: "Entrance",
								type:        "number",
								value:       this.param.entrance,
							}),
						},
					].map(object => m("div", {
						style: {
							display:    "table-row",
							whiteSpace: "nowrap",
						},
					},
						m("div", {
							className: "component-app-members-cell",
							style:     {
								padding:       "1rem",
								verticalAlign: "middle",
							},
						}, object.label),
						m("div", {
							className: "component-app-members-cell",
							style:     {
								padding: "1rem",
								width:   "100%",
							},
						}, object.input)
					))),
					m("div", {
						style: {
							display:        "flex",
							justifyContent: "space-around",
						},
					}, [
						{
							label: "OBのみ",
							value: "1",
						}, {
							label: "現役のみ",
							value: "0",
						}, {
							label: "OB/現役不問",
							value: null,
						},
					].map(object => m("label", {className: "form-group"},
						m("input", {
							checked:   object.value == this.param.ob,
							className: "radio-inline",
							name:      "ob",
							oninput:   m.withAttr("checked",
								checked => checked && update.call(this, "ob", object.value)),
							style: {margin: "0"},
							type:  "radio",
						}), m("span", {
							className: "control-label",
						}, " ", object.label))))),
				this.error && m("div", {
					className: "alert alert-danger",
					role:      "alert",
				},
					m("span", {"aria-hidden": "true"},
						m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
						" "),
					this.error),
				this.members && m("div",
					m("p", {
						className: "lead",
						style:     {color: "gray"},
					}, this.count + " 件"),
					m("div", {style: {margin: "1rem"}},
						m("a", {
							className: "btn btn-primary",
							href:      "#!member",
							onclick:   addMember.bind(this),
							role:      "button",
						},
							m("span", {"aria-hidden": "true"},
								m("span", {className: "glyphicon glyphicon-plus"}),
								" "),
							"追加")),
					m("div", {style: {margin: "1rem"}},
						 m(this.table, {
							members:     this.members,
							oncreate:    setTableState.bind(this),
						}))))),
		m("div", {
			"aria-hidden": (!promise || promise.state() != "pending").toString(),
			id:            "component-app-members-loading",
			style:         {display: "none"},
		}, "読み込み中…"),
	];
}

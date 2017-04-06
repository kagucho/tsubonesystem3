/**
	@file members_picker.js implements the feature to pick members.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/membersPicker */

import * as alert from "./alert";
import * as modal from "../modal";
import * as progress from "../progress";
import Stream from "mithril/stream";
import Table from "./table";
import client from "../client";

/**
	labelID is a string of the ID labelling members picker.
	@type !String
*/
export const labelID = "component-members-picker-title";

/**
	Internal represents the internal state.
	@private
	@extends Object
*/
class Internal {
	/**
		constructor constructs Internal.
		@returns {Undefined}
	*/
	constructor(membersStream) {
		this.clubs = new Map;
		this.count = Stream(0);
		this.members = new Map;
		this.membersStream = membersStream;
		this.table = new Table(this.tableView.bind(this));
	}

	/**
		TODO
	*/
	deinit() {
		for (const stream of [
			this.clubsStream,
			this.mappedMembersStream,
			this.officersStream,
		]) {
			stream.end(true);
		}

		if (this.loading) {
			this.loading.remove();
		}

		if (this.tableState && this.tableState.modal) {
			this.tableState.modal.remove();
		}
	}

	/**
		load loads the remote content. TODO
		@returns {Undefined}
	*/
	init() {
		this.clubsStream = client.mapClubs(
			promise => promise.then(
				clubs => {
					const newClubs = new Map;

					for (const newClub of clubs()) {
						const club = this.clubs.get(newClub.id);
						newClub.checked = !club || club.checked;
						newClubs.set(newClub.id, newClub);
					}

					this.clubs = newClubs;
					this.filter();
					m.redraw();
				}));

		this.mappedMembersStream = this.membersStream.map(
			promise => promise.then(members => {
				const newMembers = new Map;
				let count = 0;

				for (const newMember of members()) {
					const member = this.members.get(newMember.id);

					if (member && member.checked) {
						newMember.checked = true;
						count++;
					}

					newMember.clubs = [];
					newMembers.set(newMember.id, newMember);
				}

				this.count(count);
				this.members = newMembers;
				this.filter();
				m.redraw();
			}));

		this.officersStream = client.mapOfficers(promise => promise);

		client.merge(this.clubsStream, this.mappedMembersStream).map(
			promise => promise.done(
				() => {
					for (const [id, {members}] of this.clubs) {
						for (const memberID of members) {
							const member = this.members.get(memberID);
							if (member) {
								member.clubs.push(id);
							}
						}
					}

					this.filter();
					m.redraw();
				}));

		client.merge(this.officersStream, this.clubsStream, this.mappedMembersStream).map(promise => {
			this.loading = progress.add({
				"aria-describedby": "component-members-picker-loading",
				value:              0,
			});

			promise.then(
				officers => {
					const officerSet = new Set;

					for (const officer of officers()) {
						officerSet.add(officer.member);
					}

					for (const club of this.clubs.values()) {
						officerSet.add(club.chief);
					}

					for (const member of this.members.values()) {
						member.officer = officerSet.has(member.id);
					}

					this.filter();
					this.loading.remove();
					m.redraw();
				},
				error => modal.add(
					alert.closable({onclosed: this.loading.remove},
						m("span", {"aria-hidden": "true"},
							m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
							" "),
						client.error(error))),
				event => this.loading.updateValue(
					{max: event.total, value: event.loaded}));
		});
	}

	/**
		filter filters members to show.
		@returns {Undefined}
	*/
	filter() {
		for (const member of this.members.values()) {
			member.hidden = (this.ob != null && member.ob != this.ob) ||
				(this.officer != null && member.officer != this.officer) ||
				!member.clubs.some(club => this.clubs.get(club).checked);
		}
	}

	/**
		notMember reverses the checked state of a member.
		@param {!*} entry - The entry of the member.
		@returns {Undefined}
	*/
	notMember(entry) {
		entry.checked = !entry.checked;
		if (entry.checked) {
			this.count(this.count() + 1);
		} else {
			this.count(this.count() - 1);
		}
	}

	/**
		checkVisibleMembers checks the visible members.
		@returns {Undefined}
	*/
	checkVisibleMembers() {
		for (const member of this.members.values()) {
			if (!member.hidden && !member.checked) {
				member.checked = true;
				this.count(this.count() + 1);
			}
		}
	}

	/**
		uncheckVisibleMembers unchecks the visible members.
		@returns {Undefined}
	*/
	uncheckVisibleMembers() {
		for (const member of this.members.values()) {
			if (!member.hidden && member.checked) {
				member.checked = false;
				this.count(this.count() - 1);
			}
		}
	}

	/**
		checkAllMembers checks all the members.
		@returns {Undefined}
	*/
	checkAllMembers() {
		for (const member of this.members.values()) {
			member.checked = true;
		}

		this.count(this.members.size);
	}

	/**
		uncheckAllMembers unchecks all the members.
		@returns {Undefined}
	*/
	uncheckAllMembers() {
		for (const member of this.members.values()) {
			member.checked = false;
		}

		this.count(0);
	}

	/**
		notClub reverses the checked state of a club condition.
		@param {!*} entry - The entry of the club.
		@returns {Undefined}
	*/
	notClub(entry) {
		entry.checked = !entry.checked;
		this.filter();
	}

	/**
		updateOB updates ob condition.
		@param {?Boolean} value - A Boolean indcating whether the
		visible members should be OBs or not. If it is null, the
		condition will be ignored.
		@returns {Undefined}
	*/
	updateOB(value) {
		this.ob = value;
		this.filter();
	}

	/**
		updateOfficer updates officer condition.
		@param {?Boolean} value - A Boolean indicating whether the
		visible members should be officers or not. If it is null, the
		condition will be ignored.
		@returns {Undefined}
	*/
	updateOfficer(value) {
		this.officer = value;
		this.filter();
	}

	/**
		setFocusDOM sets the DOM element to be focused when the picker
		gets focused.
		@param {!external:Mithril~Node} node - The node of the element
		to be focused when the picker gets focused.
		@returns {Undefined}
	*/
	setFocusDOM(node) {
		this.focusDOM = node.dom;
	}

	/**
		setTableState sets the state of the table.
		@param {!external:Mithril~Node} node - The node of the table.
		@returns {Undefined}
	*/
	setTableState(node) {
		this.tableState = node.state;
	}

	/**
		tableView returns the view of the table.
		@param {!*} entries - The entries to draw.
		@param {!module:private/component/table~NicknameView} - A
		function which returns the view of a member.
		@returns {!external:Mithril~Children} The view of a table.
	*/
	tableView(entries, nicknameView) {
		return [
			m("thead", m("tr", {style: {backgroundColor: "#d9edf7"}},
				m("th", "選択"),
				m("th", "ニックネーム"),
				m("th", "所属部"))),
			m("tbody", Array.from(entries, ([id, value]) => m("tr", {
				"aria-hidden": value.hidden ? "true" : "false",
				key:           id,
				style:         {
					display: value.hidden ?
						"none" : "table-row",
				},
			},
				m("td",
					m("input", {
						checked: value.checked,
						onclick: this.notMember.bind(this, value),
						type:    "checkbox",
					})),
				m("td",
					nicknameView(value),
					value.officer ?
						[" ", m("span", {className: "label label-primary"}, "えらい人")] :
						null,
					value.ob ?
						[" ", m("span", {className: "label label-default"}, "OB")] :
						null),
				m("td", value.clubs.map(id => this.clubs.get(id).name).join(", "))))),
		];
	}

	/**
		view returns the view of the picker.
		@param {!module:private/component/membersPicker} state - The
		public state. TODO
		@returns {!external:Mithril~Children} The view of the picker.
	*/
	view(title) {
		function callOnchecked(event) {
			if (event.target.checked) {
				this();
			}
		}

		return m("div", {className: "modal-content"},
			m("div", {className: "modal-header"},
				m("button", {
					ariaLabel:      "閉じる",
					className:      "close",
					"data-dismiss": "modal",
					oncreate:       this.setFocusDOM.bind(this),
				}, m("span", {ariaHidden: "true"}, "×")),
				m("div", {
					className: "lead modal-title",
					id:        labelID,
				}, title)),
			m("div", {className: "modal-body"},
				m("div", {
					role:  "search",
					style: {
						backgroundColor: "#f0f8ff",
						borderRadius:    "1rem",
						margin:          "1rem",
						padding:         "1rem",
					},
				},
					m("div", {
						style: {
							borderCollapse: "collapse",
							display:        "table",
						},
					}, [
						{
							label: m("div", {className: "control-label"}, "所属部"),
							input: m("div", {
								style: {
									display:        "flex",
									justifyContent: "space-around",
								},
							}, this.clubs && Array.from(this.clubs.values(), club => m("label",
								m("input", {
									checked: club.checked,
									onclick: this.notClub.bind(this, club),
									type:    "checkbox",
								}),
								" ",
								club.name))),
						}, {
							label: m("div", {className: "control-label"}, "現役/OB"),
							input: m("div", {
								style: {
									display:        "flex",
									justifyContent: "space-around",
								},
							}, [
								{
									label: "OBのみ",
									value: true,
								}, {
									label: "現役のみ",
									value: false,
								}, {
									label: "OB/現役不問",
									value: null,
								},
							].map(object => m("label", {
								className: "form-group",
								style:     {margin: 0},
							},
								m("input", {
									checked:   this.ob == object.value,
									className: "radio-inline",
									name:      "ob",
									onclick:   callOnchecked.bind(this.updateOB.bind(this, object.value)),
									style:     {margin: "0"},
									type:      "radio",
								}), m("span", {
									className: "control-label",
								}, " ", object.label)))),
						}, {
							label: "えらい人",
							input: m("div", {
								style: {
									display:        "flex",
									justifyContent: "space-around",
								},
							}, [
								{
									label: "えらい人のみ",
									value: true,
								}, {
									label: "普通の人も",
									value: null,
								},
							].map(object => m("label", {
								className: "form-group",
								style:     {margin: "0"},
							},
								m("input", {
									checked:   this.officer == object.value,
									className: "radio-inline",
									name:      "officer",
									onclick:   callOnchecked.bind(this.updateOfficer.bind(this, object.value)),
									style:     {margin: "0"},
									type:      "radio",
								}), m("span", {
									className: "control-label",
								}, " ", object.label)))),
						},
					].map(object => m("div", {
						style: {
							border:      "solid #DDD",
							borderWidth: "0.1rem 0",
							display:     "table-row",
							whiteSpace:  "nowrap",
						},
					},
						m("div", {
							className: "component-app-mail-cell",
							style:     {
								padding:       "1rem",
								verticalAlign: "middle",
							},
						}, object.label),
						m("div", {
							className: "component-app-mail-cell",
							style:     {
								padding: "1rem",
								width:   "100%",
							},
						}, object.input))))),
				m("button", {
					className: "btn btn-default",
					onclick:   this.checkVisibleMembers.bind(this),
					style:     {margin: "1rem"},
				},
					m("span", {ariaHidden: "true"},
						m("span", {className: "glyphicon glyphicon-check"}),
						" "),
					"引っかかったやつにチェックをつける"),
				m("button", {
					className: "btn btn-default",
					onclick:   this.uncheckVisibleMembers.bind(this),
					style:     {margin: "1rem"},
				},
					m("span", {ariaHidden: "true"},
						m("span", {className: "glyphicon glyphicon-unchecked"}),
						" "),
					"引っかかったやつのチェックを外す"),
				m("button", {
					className:      "btn btn-primary",
					style:          {margin: "1rem"},
					"data-dismiss": "modal",
				},
					m("span", {ariaHidden: "true"},
						m("span", {className: "glyphicon glyphicon-ok"}),
						" "),
					"チェックしたやつらに決定"),
				this.members && m(this.table, {
					members:     this.members,
					oncreate:    this.setTableState.bind(this),
				}),
				m("div", {
					"aria-hidden": (!this.loading || this.loading.removed()).toString(),
					id:            "component-members-picker-loading",
					style:         {display: "none"},
				}, "読み込み中…")),
			m("div", {
				className: "modal-footer",
				style:     {
					display:        "flex",
					justifyContent: "space-between",
				},
			},
				m("div",
					m("button", {
						className: "btn btn-primary",
						onclick:   this.checkAllMembers.bind(this),
					},
						m("span", {ariaHidden: "true"},
							m("span", {className: "glyphicon glyphicon-check"}),
							" "),
						"全員のチェックをつける"),
					m("button", {
						className: "btn btn-danger",
						onclick:   this.uncheckAllMembers.bind(this),
					},
						m("span", {ariaHidden: "true"},
							m("span", {className: "glyphicon glyphicon-unchecked"}),
							" "),
						"全員のチェックを外す")),
				m("button", {
					className:      "btn btn-default",
					"data-dismiss": "modal",
				}, "閉じる")));
	}

	/**
		clubs is descriptions of the clubs and the conditions to filter
		members by belonging clubs.
		@member module:private/component/membersPicker~Internal#clubs
	*/

	/**
		count is the count of the checked members. TODO
		@member {!Number} module:private/component/membersPicker~Internal#count
	*/

	/**
		focusDOM is the DOM to be focused when the picker gets focused.
		@member {?external:DOM~HTMLElement}
			module:private/component/membersPicker~Internal#focusDOM
	*/

	/**
		members is descriptions, checked state, and visibilities of the
		members.
		@member module:private/component/membersPicker~Internal#members
	*/

	/**
		ob is a condition to filter members by whether they are OB or
		not. If it is true, the visible members should be OB. If it is
		false, the visible members should not be OB. If it is null,
		the conditiong will be ignored.
		@member {?Boolean} module:private/component/membersPicker~Internal#ob
	*/

	/**
		officers is a condition to filter members by whether they are
		officers or not. If it is true, the visible members should be
		OB. If it is false, the visible members should not be OB. If it
		is null, the conditiong will be ignored.
		@member {?Boolean} module:private/component/membersPicker~Internal#officer
	*/

	/**
		table is a component to draw a table of the members.
		@member {!external:Mithril~Component}
			module:private/component/membersPicker~Internal#table
	*/

	/**
		tableState is the state of the table.
		@member {?module:private/component/table~State}
			module:private/component/membersPicker~Internal#tableState
	*/
}

/**
	newPicker returns a new picker.
	@param TODO
	@returns {module:private/component/membersPicker~Picker}
*/
export default membersStream => {
	const internal = new Internal(membersStream);

	return {
		onmodalshow() {
			internal.init();
		},

		onmodalremove() {
			internal.deinit();
		},

		onmodalshown() {
			internal.focusDOM.focus();
		},

		view() {
			return internal.view(this.title);
		},

		mapCount(callback) {
			return internal.count.map(callback);
		},

		*get() {
			for (const member of internal.members.values()) {
				if (member.checked) {
					yield member;
				}
			}
		},
	};
}

/**
	Picker is the exposed interface of an instance.
	@interface module:private/component/membersPicker~Picker
	@extends external:Mithril~Component
*/

/**
	count returns the count of the checked members.
	@function module:private/component/membersPicker~Picker#count
	@returns {!Number} The count of the checked members.
*/

/**
	getPickeds returns an iterator yielding the picked members.
	@function module:private/component/membersPicker~Picker#getPickeds
	@returns {!Iterator} An iterator yielding the checked members.
*/

/**
	TODO
	@function module:private/component/membersPicker~Picker#getAll
*/

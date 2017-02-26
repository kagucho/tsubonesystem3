/**
	@file members_picker.js implements the feature to pick members.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/membersPicker */

import * as promise from "../promise";
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
	constructor() {
		this.count = 0;
		this.table = new Table(this.tableView.bind(this));
	}

	/**
		load loads the remote content.
		@returns {Undefined}
	*/
	load() {
		const loadings = [
			client.memberListroles(),
			client.clubList(),
			client.officerList(),
		];

		const initialProgress = {computable: false, loaded: 0, total: 0};
		const progresses = [initialProgress, initialProgress, initialProgress];
		const deferred = $.Deferred();

		loadings.forEach((loading, index) => loading.progress(event => {
			progresses[index] = event;
			const snapshot = progresses.slice();

			deferred.notify(Object.defineProperties({}, {
				computable: {
					get() {
						for (const progress of snapshot) {
							if (progress.computable) {
								return progress.computable;
							}
						}

						return false;
					},
				},

				loaded: {
					get: snapshot.reduce.bind(snapshot, (sum, progress) => sum + progress.loaded, 0),
				},

				total: {
					get: snapshot.reduce.bind(snapshot, (sum, progress) => sum + progress.total, 0),
				},
			}));
		}));

		promise.when(...loadings).done((members, clubs, officers) => {
			const officerSet = new Set;
			this.clubs = {};

			for (const club of clubs[0]) {
				officerSet.add(club.chief.id);
				this.clubs[club.id] = $.extend({checked: true}, club);
			}

			for (const officer of officers[0]) {
				officerSet.add(officer.member.id);
			}

			for (const member of members[0]) {
				member.officer = officerSet.has(member.id);
			}

			[this.members] = members;
			deferred.resolve();
		}, function() {
			deferred.reject(...arguments);
		});

		return promise.wrap(deferred);
	}

	/**
		registerLoading registers a loading.
		@param {!external:jQuery~Promise} loading - A promise describing
		the loading.
		@returns {!external:jQuery~Promise} A filtered promise.
	*/
	registerLoading(loading) {
		return loading.then(
			submission => submission &&
				this.load().then(() => submission));
	}

	/**
		filter filters members to show.
		@returns {Undefined}
	*/
	filter() {
		for (const member of this.members) {
			member.hidden = true;

			if ((this.ob == null || member.ob == this.ob) &&
				(this.officer == null && member.officer == this.officer)) {
				for (const club of member.clubs) {
					if (this.clubs[club].checked) {
						member.hidden = false;
						break;
					}
				}
			}
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
			this.count++;
		} else {
			this.count--;
		}
	}

	/**
		checkVisibleMembers checks the visible members.
		@returns {Undefined}
	*/
	checkVisibleMembers() {
		for (const member of this.members) {
			if (!member.hidden && !member.checked) {
				member.checked = true;
				this.count++;
			}
		}
	}

	/**
		uncheckVisibleMembers unchecks the visible members.
		@returns {Undefined}
	*/
	uncheckVisibleMembers() {
		for (const member of this.members) {
			if (!member.hidden && member.checked) {
				member.checked = false;
				this.count--;
			}
		}
	}

	/**
		checkAllMembers checks all the members.
		@returns {Undefined}
	*/
	checkAllMembers() {
		for (const member of this.members) {
			member.checked = true;
		}

		this.count = this.members.length;
	}

	/**
		uncheckAllMembers unchecks all the members.
		@returns {Undefined}
	*/
	uncheckAllMembers() {
		for (const member of this.members) {
			member.checked = false;
		}

		this.count = 0;
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
		@param {!module:private/components/table~NicknameView} - A
		function which returns the view of a member.
		@returns {!external:Mithril~Children} The view of a table.
	*/
	tableView(entries, nicknameView) {
		return [
			m("thead", m("tr", {style: {backgroundColor: "#d9edf7"}},
				m("th", "選択"),
				m("th", "ニックネーム"),
				m("th", "所属部")
			)), m("tbody", entries.map(entry => m("tr", {
				"aria-hidden": entry.hidden ? "true" : "false",
				key:           entry.id,
				style:         {
					display: entry.hidden ?
						"none" : "table-row",
				},
			},
				m("td",
					m("input", {
						checked: entry.checked,
						onclick: this.notMember.bind(this, entry),
						type:    "checkbox",
					})
				),
				m("td",
					nicknameView(entry),
					entry.officer ?
						[" ", m("span", {className: "label label-primary"}, "役員")] :
						null,
					entry.ob ?
						[" ", m("span", {className: "label label-default"}, "OB")] :
						null
				),
				m("td", entry.clubs.map(id => this.clubs[id].name).join(", "))
			))),
		];
	}

	/**
		view returns the view of the picker.
		@param {!module:private/components/membersPicker} state - The
		public state.
		@returns {!external:Mithril~Children} The view of the picker.
	*/
	view(state) {
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
				}, state.title)
			), m("div", {className: "modal-body"},
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
							}, this.clubs && $.map(this.clubs, club => m("label",
								m("input", {
									checked: club.checked,
									onclick: this.notClub.bind(this, club),
									type:    "checkbox",
								}),
								" ",
								club.name
							))),
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
								}, " ", object.label)
							))),
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
								}, " ", object.label)
							))),
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
						}, object.input)
					)))
				), m("button", {
					className: "btn btn-default",
					onclick:   this.checkVisibleMembers.bind(this),
					style:     {margin: "1rem"},
				},
					m("span", {ariaHidden: "true"},
						m("span", {className: "glyphicon glyphicon-check"}),
						" "
					), "引っかかったやつにチェックをつける"
				), m("button", {
					className: "btn btn-default",
					onclick:   this.uncheckVisibleMembers.bind(this),
					style:     {margin: "1rem"},
				},
					m("span", {ariaHidden: "true"},
						m("span", {className: "glyphicon glyphicon-unchecked"}),
						" "
					), "引っかかったやつのチェックを外す"
				), m("button", {
					className:      "btn btn-primary",
					style:          {margin: "1rem"},
					"data-dismiss": "modal",
				},
					m("span", {ariaHidden: "true"},
						m("span", {className: "glyphicon glyphicon-ok"}),
						" "
					), "チェックしたやつらに決定"
				), this.members && m(this.table, {
					members:     this.members,
					oncreate:    this.setTableState.bind(this),
					onloadstart:
						loading => state.onloadstart(
							this.registerLoading(loading)),
				})
			), m("div", {
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
							" "
						), "全員のチェックをつける"
					), m("button", {
						className: "btn btn-danger",
						onclick:   this.uncheckAllMembers.bind(this),
					},
						m("span", {ariaHidden: "true"},
							m("span", {className: "glyphicon glyphicon-unchecked"}),
							" "
						), "全員のチェックを外す"
					)
				), m("button", {
					className:      "btn btn-default",
					"data-dismiss": "modal",
				}, "閉じる")
			)
		);
	}

	/**
		clubs is descriptions of the clubs and the conditions to filter
		members by belonging clubs.
		@member module:private/components/membersPicker~Internal#clubs
	*/

	/**
		count is the count of the checked members.
		@member {!Number} module:private/components/membersPicker~Internal#count
	*/

	/**
		focusDOM is the DOM to be focused when the picker gets focused.
		@member {?external:DOM~HTMLElement}
			module:private/components/membersPicker~Internal#focusDOM
	*/

	/**
		members is descriptions, checked state, and visibilities of the
		members.
		@member module:private/components/membersPicker~Internal#members
	*/

	/**
		ob is a condition to filter members by whether they are OB or
		not. If it is true, the visible members should be OB. If it is
		false, the visible members should not be OB. If it is null,
		the conditiong will be ignored.
		@member {?Boolean} module:private/components/membersPicker~Internal#ob
	*/

	/**
		officers is a condition to filter members by whether they are
		officers or not. If it is true, the visible members should be
		OB. If it is false, the visible members should not be OB. If it
		is null, the conditiong will be ignored.
		@member {?Boolean} module:private/components/membersPicker~Internal#officer
	*/

	/**
		table is a component to draw a table of the members.
		@member {!external:Mithril~Component}
			module:private/components/membersPicker~Internal#table
	*/

	/**
		tableState is the state of the table.
		@member {?module:private/components/table~State}
			module:private/components/membersPicker~Internal#tableState
	*/
}

/**
	newPicker returns a new picker.
	@returns {module:private/components/membersPicker~Picker}
*/
export function newPicker() {
	const internal = new Internal;

	return {
		onmodalremove() {
			if (internal.tableState && internal.tableState.modal) {
				internal.tableState.modal.remove();
			}
		},

		onmodalshown() {
			internal.focusDOM.focus();
		},

		view() {
			return internal.view(this);
		},

		count() {
			return internal.count;
		},

		*get() {
			for (const member of internal.members) {
				if (member.checked) {
					yield member;
				}
			}
		},

		load() {
			return internal.load();
		},
	};
}

/**
	Picker is the exposed interface of an instance.
	@interface module:private/components/membersPicker~Picker
	@extends external:Mithril~Component
*/

/**
	count returns the count of the checked members.
	@function module:private/components/membersPicker~Picker#count
	@returns {!Number} The count of the checked members.
*/

/**
	get returns an iterator yielding the checked members.
	@function module:private/components/membersPicker~Picker#get
	@returns {!Iterator} An iterator yielding the checked members.
*/

/**
	load loads the remote content.
	@function module:private/components/membersPicker~Picker#load
	@returns {Undefined}
*/

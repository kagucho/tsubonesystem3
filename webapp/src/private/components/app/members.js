/**
	@file members.js implements the feature to show members.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/members */

/**
	module:private/components/app/members is a component to show members.
	@name module:private/components/app/members
	@type external:Mithril~Component
*/

import * as client from "../../client";
import * as container from "../container";
import * as modal from "../member/modal";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";
import large from "../../large";
import {members} from "../table";

export class controller {
	constructor() {
		this.param = m.route.param();
		this.progress = new ProgressSum;
		this.load();
	}

	addMember() {
		if (large()) {
			this.adding = true;

			return false;
		}
	}

	endAddingMember() {
		delete this.adding;
	}

	load() {
		this.progress.add(client.memberList().then(got => {
			this.members = got;
		}, xhr => {
			this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
		}));
	}

	update(key, value) {
		if (value) {
			this.param[key] = value;
		} else {
			delete this.param[key];
		}

		let url = "#!members";
		if (!$.isEmptyObject(this.param)) {
			url += "?" + m.route.buildQueryString(this.param);
		}

		history.pushState(null, null, url);
	}
}

export function view(control) {
	const content = [
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
							control.update.bind(control, "nickname")),
						placeholder: "Nickname",
						type:        "search",
						value:       control.param.nickname || "",
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
							control.update.bind(control, "realname")),
						placeholder: "Name",
						type:        "search",
						value:       control.param.realname || "",
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
							control.update.bind(control, "entrance")),
						placeholder: "Entrance",
						type:        "number",
						value:       control.param.entrance,
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
				},
				{
					label: "現役のみ",
					value: "0",
				},
				{
					label: "OB/現役不問",
					value: null,
				},
			].map(object => m("label", {
				className: "form-group",
			},
				m("input", {
					checked:   object.value == control.param.ob,
					className: "radio-inline",
					name:      "ob",
					oninput:   m.withAttr("checked", function(checked) {
						if (checked) {
							this.control.update("ob", this.value);
						}
					}.bind({control, value: object.value})),
					style: {margin: "0"},
					type:  "radio",
				}), m("span", {
					className: "control-label",
				}, " ", object.label)
			)))
		),
	];

	const onloadstart = (function(promise) {
		this.progress.add(promise.done(
			submission => submission && this.load()));
	}).bind(control);

	if (control.error) {
		content.push(
			m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
			m("span", {ariaHidden: "true"},
				m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
				" "
			), control.error
		));
	}

	if (control.members) {
		const showingMembers = new Array(control.members.length);
		let count = 0;

		for (const member of control.members) {
			if ((!control.param.entrance ||
					member.entrance == control.param.entrance) &&
				(!control.param.nickname ||
					member.nickname.includes(control.param.nickname)) &&
				(!control.param.ob ||
					control.param.ob == "1" && member.ob ||
					control.param.ob == "0" && !member.ob) &&
				(!control.param.realname ||
					member.realname.includes(control.param.realname))) {
				showingMembers[count] = member;
				count++;
			}
		}

		content.push(m("div",
			m("p", {
				className: "lead",
				style:     {
					clear: "both",
					color: "gray",
				},
			}, count + " 件"),
			m("div", {style: {margin: "1rem"}},
				m("a", {
					className: "btn btn-primary",
					href:      "#!member",
					onclick:   control.addMember.bind(control),
				}, "追加")
			),
			m("div", {
				className: "table-responsive",
				style:     {margin: "1rem"},
			},
				 m(members, {
					members: showingMembers,
					onloadstart,
				})
			)
		));
	}

	return [
		m(progress, control.progress.html()),
		m(container, m("div", {className: "container"},
			m("h1", "Members"),
			m("div", content)
		)),
		control.adding && m(modal, {
			id:       null,
			onhidden: control.endAddingMember.bind(control),
			onloadstart,
		}),
	];
}

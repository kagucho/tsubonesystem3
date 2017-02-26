/**
	@file party.js implements party component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

import * as alert from "../alert";
import * as container from "../container";
import * as membersPicker from "../members_picker";
import * as modal from "../../modal";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";
import client from "../../client";

/**
	TODO
*/
function defaultInviteds() {
	if (this.defaultInviteds !== undefined) {
		this.defaultInviteds =
			this.membersPicker.count() + "人のフレンズ";
	}
}

/**
	TODO
*/
function updateInviteds(target) {
	if (target.value == "") {
		this.defaultInviteds = null;
		defaultInviteds.call(this);
	} else if (target.value != this.defaultInviteds) {
		delete this.defaultInviteds;
	}
}

/**
	TODO
*/
function registerLoading(promise) {
	this.progress.add(promise.catch(xhr => {
		openError.call(this, alert.closable,
			client.error(xhr) || "どうしようもないエラーです");
	}));
}

/**
	TODO
*/
function submit(target) {
	let valid = true;
	for (const key in this.validities) {
		if (!this.validities[key].reportValidity(this, target[key])) {
			valid = false;
		}
	}

	if (!valid) {
		return;
	}

	const param = {};

	/* eslint-disable camelcase */
	param.name = target.name.value;
	param.place = target.place.value;
	param.inviteds_name = target.inviteds.value;
	param.details = target.details.value;
	/* eslint-enable camelcase */

	const datetime = $(target.datetime).data("daterangepicker");

	param.start = datetime.startDate.unix();
	param.end = datetime.endDate.unix();
	param.due = $(target.due).data("daterangepicker").startDate.unix();

	const members = this.membersPicker.get();

	const iteration = members.next();
	param.inviteds = iteration.value.id;

	for (const member of members) {
		param.inviteds = [param.inviteds, member.id].join(" ");
	}

	const inprogress = openInprogress.call(this, "送信しています…");

	this.progress.add(client.partyCreate(param).then(response => {
		inprogress.remove();

		if (response.error == "mail_failure") {
			openError.call(this, alert.leavable,
				"メールの送信に失敗しました");
		} else {
			openOK.call(this, alert.leavable, "送信しました。");
		}
	}, xhr => {
		inprogress.remove();

		openError.call(this, alert.closable,
			client.error(xhr) || "どうしようもないエラーです");
	}));
}

/**
	TODO
*/
function openMembersPicker() {
	if (!this.modalMembersPicker || this.modalMembersPicker.removed()) {
		this.modalMembersPicker = modal.unshift(
			{"aria-labelledby": membersPicker.labelID},
			this.membersPicker);
	}
}

/**
	TODO
*/
function openError(specifiedAlert, ...children) {
	modal.unshift(specifiedAlert(
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
			" "
		), ...children
	));
}

/**
	TODO
*/
function openOK(specifiedAlert, ...children) {
	modal.unshift(specifiedAlert(
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-ok"}),
			" "
		), ...children
	));
}

/**
	TODO
*/
function openInprogress() {
	return modal.unshift({backdrop: "static"}, alert.inprogress(...arguments));
}

/**
	TODO
*/
class Validity {
	/**
		TODO
	*/
	constructor() {
		this.validationMessage = null;
		this.view = this.view.bind(this);
	}

	/**
		TODO
	*/
	updateValidationMessage(state, target) {
		if (!this.validationMessage) {
			state.invalids++;
		}

		this.validationMessage = target.validationMessage;
	}

	/**
		TODO
	*/
	checkValidity(state, target) {
		const checked = target.checkValidity();

		if (checked && this.validationMessage) {
			this.validationMessage = null;
			state.invalids--;
		}

		return checked;
	}

	/**
		TODO
	*/
	dismissValidationMessage(state, target) {
		if (this.validationMessage) {
			this.checkValidity(state, target);
		}
	}

	/**
		TODO
	*/
	reportValidity(state, target) {
		const reported = target.reportValidity();

		if (reported && this.validationMessage) {
			this.validationMessage = null;
			state.invalids--;
		}

		return reported;
	}

	/**
		TODO
	*/
	view() {
		return m("div", {
			"aria-hidden": (!this.validationMessage).toString(),
			style:         {
				margin:     "1rem",
				minWidth:   "32ch",
				visibility: this.validationMessage ?
					"visible" : "hidden",
			},
		},
			m("div", {
				className: "alert alert-danger",
				role:      "alert",
				style:     {
					display: "inline-block",
					margin:  "0",
					padding: ".5rem",
				},
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), this.validationMessage
			)
		);
	}
}

/**
	TODO
*/
class CustomValidity extends Validity {
	checkValidity(state, target) {
		this.constructor.setCustomValidity(state, target);

		return super.checkValidity(state, target);
	}

	reportValidity(state, target) {
		this.constructor.setCustomValidity(state, target);

		return super.reportValidity(state, target);
	}
}

/**
	TODO
*/
class NameValidity extends CustomValidity {
	static setCustomValidity(state, target) {
		if (state.parties.has) {
			target.setCustomValidity(state.parties.has(target.value) ?
				"その題目、ダブってるよ。" : "");
		}
	}
}

/**
	TODO
*/
class DatetimeValidity extends CustomValidity {
	static setCustomValidity(state, target) {
		const datetime = $(target).data("daterangepicker");

		if (datetime) {
			target.setCustomValidity(moment().add(1, "hour").isAfter(datetime.startDate) ?
				"1時間以内に開催されるパーティーは作成できません" : "");
		}
	}

	checkValidity(state, target) {
		state.validities.due.dismissValidationMessage(state, target.form.due);

		return super.checkValidity(state, target);
	}

	dismissValidationMessage(state, target) {
		state.validities.due.dismissValidationMessage(state, target.form.due);
		super.dismissValidationMessage(state, target);
	}
}

/**
	TODO
*/
class DueValidity extends CustomValidity {
	static setCustomValidity(state, target) {
		const datetime = $(target.form.datetime).data("daterangepicker");
		const due = $(target).data("daterangepicker");

		target.setCustomValidity(due ?
			(datetime && due.startDate.isAfter(datetime.endDate) ?
				"出欠締め切りは終了前にしてください" :
				(moment().add(1, "hour").isAfter(due.startDate) ?
					"出欠締め切りは1時間以上あとにしてください" :
					""
				)
			) : "");
	}
}

export function oninit() {
	this.parties = client.partyListnames().then(names => {
		this.parties = new Set(names);
	});

	this.membersPicker = membersPicker.newPicker();
	this.membersPicker.onloadstart = registerLoading.bind(this);
	this.membersPicker.title = "誰を誘う?";
	this.membersPickerLoading = this.membersPicker.load();

	this.progress = new ProgressSum;
	registerLoading.call(this, this.parties);
	registerLoading.call(this, this.membersPickerLoading);

	this.defaultInviteds = "";

	this.invalids = 0;
	this.messages = [];
	this.pickingMembers = false;
	this.validities = {
		name:     new NameValidity,
		datetime: new DatetimeValidity,
		place:    new Validity,
		inviteds: new Validity,
		due:      new DueValidity,
		details:  new Validity,
	};
}

export function onbeforeremove() {
	if (this.modalMembersPicker) {
		this.modalMembersPicker.remove();
	}
}

export function view() {
	const now = moment();
	const nextHour = now.clone().add(1, "hour");
	const tomorrow = now.clone().startOf("day").add(1, "day");

	const membersPickerBarrier = {
		pending:  "読み込み中です",
		rejected: "読み込みに失敗したため使用できません",
	}[this.membersPickerLoading.state()] || "";

	const submissionBarrier = (this.parties.state && {
		pending:  "読み込み中です",
		rejected: "読み込みに失敗したため使用できません",
	}[this.parties.state()]) || this.invalids ?
		"無効な項目があります" :
		(this.membersPicker.count() > 0 ? "" : "誰か選んでください");

	function handleValidity(inputKey, functionKey, event) {
		this.validities[inputKey][functionKey](this, event.target);
	}

	const bindValidity = handleValidity.bind.bind(handleValidity, this);

	defaultInviteds.call(this);

	return [
		m(progress, this.progress.html()),
		m(container, m("div", {className: "container"},
			m("h1", "Party"),
			"新しいパーティーを作成します。",
			m("form", {
				display: "flex",

				onsubmit: event => {
					submit.call(this, event.target);

					return false;
				},
			},
				m("label", {className: "center-block form-group"},
					m("div", {className: "control-label"},
						"題目"
					), m("div", {
						style: {
							display:  "flex",
							flexWrap: "wrap",
						},
					},
						m("input", {
							className:   "form-control",
							maxlength:   "63",
							name:        "name",
							onchange:    bindValidity("name", "checkValidity"),
							oninput:     bindValidity("name", "dismissValidationMessage"),
							oninvalid:   bindValidity("name", "updateValidationMessage"),
							placeholder: "Title",
							style:       {
								flex:       "1",
								fontWeight: "400",
								margin:     "1rem",
								maxWidth:   "63ch",
							},
						}), m(this.validities.name)
					)
				), m("label", {className: "center-block form-group"},
					m("div", {className: "control-label"}, "時刻"),
					m("div", {
						style: {
							display:  "flex",
							flexWrap: "wrap",
						},
					},
						m("input", {
							className: "form-control",
							name:      "datetime",

							oncreate: node => {
								$(node.dom).daterangepicker({
									locale: {
										applyLabel:  "決定",
										cancelLabel: "やっぱやめた",
										format:      "llll",
									},

									minDate:    nextHour,
									startDate:  tomorrow,
									endDate:    tomorrow,
									timePicker: true,
								});
							},

							onchange:    bindValidity("datetime", "checkValidity"),
							oninput:     bindValidity("datetime", "dismissValidationMessage"),
							oninvalid:   bindValidity("datetime", "updateValidationMessage"),
							placeholder: "Date and time",
							style:       {
								flex:       "1",
								fontWeight: "400",
								margin:     "1rem",
								maxWidth:   "63ch",
							},
						}), m(this.validities.datetime)
					)
				), m("label", {className: "center-block form-group"},
					m("div", {className: "control-label"},
						"開催場所"
					), m("div", {
						style: {
							display:  "flex",
							flexWrap: "wrap",
						},
					},
						m("input", {
							className:   "form-control",
							maxlength:   "63",
							name:        "place",
							onchange:    bindValidity("place", "checkValidity"),
							oninput:     bindValidity("place", "dismissValidationMessage"),
							oninvalid:   bindValidity("place", "updateValidationMessage"),
							placeholder: "Place",
							style:       {
								flex:       "1",
								fontWeight: "400",
								margin:     "1rem",
								maxWidth:   "63ch",
							},
						}), m(this.validities.place)
					)
				), m("label", {className: "center-block form-group"},
					m("div", {className: "control-label"},
						"招待対象者"
					), m("div", {
						style: {
							display:    "flex",
							flexWrap:   "wrap",
							fontWeight: "400",
							width:      "100%",
						},
					},
						m("input", {
							className: "form-control",
							maxlength: "63",
							name:      "inviteds",
							onchange:  function(checkValidity, event) {
								updateInviteds.call(this, event.target);
								checkValidity(event);
							}.bind(this, bindValidity("inviteds", "checkValidity")),
							oninput:     bindValidity("inviteds", "dismissValidationMessage"),
							oninvalid:   bindValidity("inviteds", "updateValidationMessage"),
							placeholder: "Subjects",
							style:       {
								flex:     "1",
								margin:   "1rem",
								maxWidth: "63ch",
							},
							value: this.defaultInviteds,
						}), m("button", {
							className: "btn btn-default",
							disabled:  Boolean(membersPickerBarrier),
							onclick:   openMembersPicker.bind(this),
							style:     {margin: "1rem"},
							title:     membersPickerBarrier,
							type:      "button",
						},
							m("span", {"aria-hidden": "true"},
								m("span", {className: "glyphicon glyphicon-check"}),
								" "
							), "変更する"
						), m(this.validities.inviteds)
					)
				), m("label", {className: "center-block form-group"},
					m("div", {className: "control-label"},
						"出欠締め切り時刻"
					), m("div", {
						style: {
							display:  "flex",
							flexWrap: "wrap",
						},
					},
						m("input", {
							className: "form-control",
							name:      "due",

							oncreate: node => {
								$(node.dom).daterangepicker({
									locale: {
										applyLabel:  "決定",
										cancelLabel: "やっぱやめた",
										format:      "llll",
									},

									minDate:          nextHour,
									startDate:        tomorrow,
									singleDatePicker: true,
									timePicker:       true,
								});
							},

							onchange:    bindValidity("due", "checkValidity"),
							oninput:     bindValidity("due", "dismissValidationMessage"),
							oninvalid:   bindValidity("due", "updateValidationMessage"),
							placeholder: "Due",
							style:       {
								fontWeight: "400",
								margin:     "1rem",
								maxWidth:   "63ch",
							},
						}), m(this.validities.due)
					)
				), m("label", {
					className: "form-group",
					style:     {
						display:       "flex",
						flex:          "1",
						flexDirection: "column",
					},
				},
					m("div", {style: {display: "flex"}},
						m("div", {className: "control-label"}, "詳細"),
						m(this.validities.details)
					), m("textarea", {
						className:   "form-control",
						maxlength:   "8192",
						name:        "details",
						onchange:    bindValidity("details", "checkValidity"),
						oninput:     bindValidity("details", "dismissValidationMessage"),
						oninvalid:   bindValidity("details", "updateValidationMessage"),
						placeholder: "Details",
						style:       {
							flex:       "1",
							fontWeight: "400",
							margin:     "1rem",
						},
					})
				), m("button", {
					className: "btn btn-block btn-primary",
					disabled:  Boolean(submissionBarrier),
					title:     submissionBarrier,
				}, "送信")
			)
		)),
	];
}

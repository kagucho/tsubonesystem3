/**
	@file primitive.js provides the primitive elements of the UI to show or
	create email.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/mail/primitive */

import * as alert from "../alert";
import * as container from "../container";
import * as modal from "../../modal";
import * as picker from "../members_picker";
import * as progress from "../../progress";
import * as table from "../table";
import client from "../../client";

/**
	openError opens an error dialog.
	@private
	@function
	@param {!module:private/modal~Component} specifiedAlert - A component
	to draw an alert.
	@param {...?external:Mithril~Children} children - An error message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openError = (specifiedAlert, ...children) => modal.add(specifiedAlert(
	m("span", {"aria-hidden": "true"},
		m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
		" "
	), ...children
));

/**
	openOK opens a dialog showing a successful result.
	@private
	@function
	@param {!module:private/modal~Component} specifiedAlert - A component
	to draw an alert.
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openOK = (specifiedAlert, ...children) => modal.add(specifiedAlert(
	m("span", {"aria-hidden": "true"},
		m("span", {className: "glyphicon glyphicon-ok"}),
		" "
	), ...children
));

/**
	openBusy opens a dialog showing a busy state.
	@private
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openBusy() {
	return modal.add({backdrop: "static"}, alert.busy(...arguments));
}

/**
	openPicker opens a dialog to pick members.
	@private
	@param TODO
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openPicker =
	picker => modal.add({"aria-labelledby": picker.labelID}, picker);

class Internal {
	constructor(external, subject) {
		this.external = external;
		this.subject = subject;
		this.to = $.noop;
	}

	/**
		defaultTo defaults to if necessary.
		@private
		@returns {Undefined}
	*/
	defaultTo() {
		this.to = this.picker.mapCount(count => {
			m.redraw();
			return count + "人のフレンズ";
		});
	}

	start() {
		if (this.subject) {
			this.streams = [
				client.mapMail(this.subject,
					promise => promise.done(
						mail => this.mail = mail)),

				client.mapMembers(promise => promise),
			];

			this.loading = client.merge(...this.streams).map(
				promise => promise.done((mail, members) => {
					this.members = [];

					for (const member of members()) {
						if (mail.recipients.has(member.id)) {
							this.members.push(member);
						}

						if (member.id == mail.from) {
							this.from = member;
						}
					}
				}));
		} else {
			this.streams = [
				client.mapMemberMails(
					promise => promise.then(
						mails => new Set((function *() {
							for (const member of mails()) {
								if (member.mail) {
									yield member.id;
								}
							}
						})()))),

				client.mapMembers(promise => promise),
			];

			this.picker = picker.default(client.merge(...this.streams).map(
				promise => promise.then(
					Function.prototype.bind.bind(function *(mails, roles) {
						for (const member of roles()) {
							if (mails.has(member.id)) {
								yield member;
							}
						}
					}, undefined))));

			this.loading = client.mapMails(
				promise => promise.done(
					mails => this.subjects = new Set((function *() {
						for (const mail of mails()) {
							yield mail.subject;
						}
					})())));

			this.streams.push(this.loading);
			this.count = this.picker.mapCount(count => count);
			this.defaultTo();
			this.picker.title = "誰に送る?";
		}

		this.loading.map(promise => {
			const loadingProgress = progress.add({
				"aria-describedby": "component-app-mail-loading",
				value:              0,
			});

			promise.then((mail, members) => {
				loadingProgress.remove();
				m.redraw();
			}, error => {
				loadingProgress.updateARIA({"aria-describedby": alert.bodyID});
				m.redraw();

				openError(alert.closable.bind({onclosed: loadingProgress.remove}),
					client.error(error));
			}, event => loadingProgress.updateValue(
				{max: event.total, value: event.loaded}));

			m.redraw();
		});
	}

	end() {
		if (this.modalPicker) {
			this.modalPicker.remove();
		}

		for (const stream of this.streams) {
			stream.end(true);
		}
	}

	/**
		TODO
	*/
	openPicker() {
		this.modalPicker = openPicker(this.picker);
	}

	/**
		submit submits the mail.
		@param {!external:DOM~HTMLFormElement} target - A form representing
		properties of the mail.
		@returns {Undefined}
	*/
	submit(target) {
		if (!target.reportValidity()) {
			return;
		}

		const param = {body: target.body.value, to: target.to.value};
		const pickeds = this.picker.get();

		const iteration = pickeds.next();
		param.recipients = iteration.value.id;

		for (const picked of pickeds) {
			param.recipients = [param.recipients, picked.id].join(" ");
		}

		const submissionModal = openBusy("送信しています…");
		const submissionProgress = progress.add(
			{"aria-describedby": alert.bodyID, value: 0});

		client.createMail(target.subject.value, param).then(response => {
			submissionModal.remove();

			if (response.error == "mail_failure") {
				openError(
					alert.closable.bind(
						{onclosed: submissionProgress.remove}),
					"メールの送信に失敗しました");
			} else {
				openOK(this.external.onemptied ?
						alert.closable.bind({
							onclosed() {
								submissionProgress.remove();
								this.external.onemptied();
							},
						}) :
						alert.leavable.bind({
							onclosed() {
								submissionProgress.remove();
							},
						}),
					"送信しました");
			}
		}, error => {
			submissionModal.remove();

			openError(
				alert.closable.bind(
					{onclosed: submissionProgress.remove}),
				client.error(error));
		}, event => {
			submissionProgress.updateValue(
				{max: event.total, value: event.loaded});

			event.redraw = false;
		});
	}

	/**
		TODO
	*/
	updateTo(value) {
		if (!value) {
			this.defaultTo();
		} else if (value != this.to()) {
			this.to.end(true);
			this.to();
		}
	}

	checkSubjectValidity(target) {
		if (this.subjects) {
			target.setCustomValidity(this.subjects.has(target.value) ?
				"その題名、ダブってるよ" : "");
		}
	}

	dismissSubjectValidity(target) {
		if (target.validationMessage && !this.subject.has(target.value)) {
			tagret.setCustomValidity(target);
		}
	}

	/**
		TODO
	*/
	bodyView() {
		const label = this.subject ? "div" : "label";
		const loadingPromise = this.loading && this.loading();

		return [
			m("div", {
				style: {
					display:       "flex",
					flexDirection: "column",
					minHeight:     "100%",
				},
			},
				this.mail && [
					m("div", {className: "form-group"},
						m("div", {
							className: "control-label",
							style:     {fontWeight: "bold"},
						}, "Date"),
						moment.unix(this.mail.date).format("lll")),
					m("div", {className: "form-group"},
						m("div", {
							className: "control-label",
							style:     {fontWeight: "bold"},
						}, "From"),
						this.from && m("a", {href: "#!member?id=" + this.from.id},
							this.from.nickname)),
				],
				m("div", {className: "form-group"},
					m("div", {
						className: "control-label",
						style:     {fontWeight: "bold"},
					}, "To"),
					this.subject ?
						this.mail && m("div", {
							id:        "component-mail-to-group",
							className: "panel-group",
						},
							m("div", {className: "panel panel-default"},
								m("div", {
									className: "panel-heading",
									id:        "component-mail-to-heading",
								},
									m("div", {className: "panel-title"},
										m("button", {
											"aria-controls": "component-mail-to-collapse",
											"data-parent":   "#component-mail-to-group",
											"data-toggle":   "collapse",
											oncreate(node) {
												node.dom.setAttribute("aria-expanded", "false");
											},
										}, this.mail.to))),
								m("div", {
									"aria-labelledby": "component-mail-to-heading",
									id:                "component-mail-to-collapse",
									oncreate(node) {
										node.dom.className = "panel-collapse collapse";
									},
								}, m("div", {className: "panel-body"},
									m(table.members, {members: this.members}))))) :
						m("div", {
							style: {
								display:  "flex",
								flexWrap: "wrap",
								width:    "100%",
							},
						},
							m("input", {
								className:   "form-control",
								maxlen:      "63",
								name:        "to",
								onblur:      m.withAttr("value", this.updateTo.bind(this)),
								placeholder: "To",
								required:    true,
								style:       {flex: "1"},
								value:       this.to(),
							}),
							m("button", {
								className: "btn btn-default",
								onclick:   this.openPicker.bind(this),
								oncreate:  node => node.dom.focus(),
								type:  "button",
							},
								m("span", {"aria-hidden": "true"},
									m("span", {className: "glyphicon glyphicon-check"}),
									" "),
								"変更する"))),
				m(label, {
					className: "form-group",
					style:     {display: "block"},
				},
					m("div", {
						className: "control-label",
						style:     {fontWeight: "bold"},
					}, "Subject"),
					this.subject || m("input", {
						className:   "form-control",
						maxlen:      "63",
						name:        "subject",
						placeholder: "Subject",
						required:    true,
						style:       {fontWeight: "normal"},
					})),
				m(label, {
					className: "form-group",
					style:     {
						display:       "flex",
						flexDirection: "column",
						flex:          "1",
					},
				},
					m("div", {
						className: "control-label",
						style:     {fontWeight: "bold"},
					}, "Body"),
					this.subject ?
						m("pre", this.mail && this.mail.body) :
						m("textarea", {
							/*
								RFC 5322 - Internet Message Format
								2.1.1.  Line Length Limits
								https://tools.ietf.org/html/rfc5322#section-2.1.1
								> Each line of characters MUST be no more than
								> 998 characters, and SHOULD be no more than 78 characters, excluding
								> the CRLF.
							*/
							cols: "78",

							className:   "form-control",
							maxlen:      "8192",
							name:        "body",
							oninput:     event => this.dismissSubjectValidity(event.target),
							onblur:      event => this.checkSubjectValidity(event.target),
							placeholder: "Body",
							required:    true,
							style:       {
								flex:       "1",
								fontWeight: "normal",
							},
						}))),
			m("div", {
				"aria-hidden": (!loadingPromise || loadingPromise.state() != "pending").toString(),
				id:            "component-app-mail-loading",
				style:         {display: "none"},
			}, "読み込み中…"),
		];
	}

	buttonView() {
		if (this.subject) {
			return null;
		}

		const loadingPromise = this.loading && this.loading();
		const barrier = {
			pending:  "読み込み中です",
			rejected: "読み込みに失敗したため使用できません",
		}[loadingPromise ? loadingPromise.state() : "pending"] ||
			(this.count() <= 0 ? "送る相手を選んでください" : "");

		return m("button", {
			className: "btn btn-primary",
			disabled:  Boolean(barrier),
			onclick:   function(event) {
				this(event.target.form);
				return false;
			}.bind(this.submit.bind(this)),
			title:     barrier,
		},
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-send"}),
				" "),
			"送信");
	}
}

/**
	TODO
*/
export default function(subject) {
	const external = {};
	const internal = new Internal(external, subject);

	external.start = () => internal.start();
	external.end = () => internal.end();

	external.body = {
		view() {
			return internal.bodyView();
		},
	};

	external.button = {
		view() {
			return internal.buttonView();
		},
	};

	return external;
}

/**
	TODO
*/
export const title = "Mail";

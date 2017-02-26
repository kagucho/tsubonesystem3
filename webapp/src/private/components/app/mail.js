/**
	@file mail.js implements mail component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/mail */

/**
	module:private/components/app/mail is a component to mail.
	@name module:private/components/app/mail
	@type !external:Mithril~Component
*/

import * as alert from "../alert";
import * as container from "../container";
import * as modal from "../../modal";
import * as picker from "../members_picker";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";
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
const openError = (specifiedAlert, ...children) => modal.unshift(specifiedAlert(
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
const openOK = (specifiedAlert, ...children) => modal.unshift(specifiedAlert(
	m("span", {"aria-hidden": "true"},
		m("span", {className: "glyphicon glyphicon-ok"}),
		" "
	), ...children
));

/**
	openInprogress opens a dialog showing progress.
	@private
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openInprogress() {
	return modal.unshift({backdrop: "static"}, alert.inprogress(...arguments));
}

/**
	openPicker opens a dialog to pick members.
	@private
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openPicker() {
	if (!this.modalPicker || this.modalPicker.removed()) {
		this.modalPicker = modal.unshift(
			{"aria-labelledby": picker.labelID},
			this.picker);
	}

	return this.modalPicker;
}

/**
	defaultTo defaults to if necessary.
	@private
	@this module:private/components/app/mail
	@returns {Undefined}
*/
function defaultTo() {
	if (this.defaultTo !== undefined) {
		this.defaultTo = this.picker.count() + "人のフレンズ";
	}
}

/**
	updateTo updates to.
	@private
	@this module:private/components/app/mail
	@param {?String} value - The value of to. If it is empty, to will have
	the default value.
	@returns {Undefined}
*/
function updateTo(value) {
	if (!value) {
		this.defaultTo = null;
		defaultTo.call(this);
	} else if (value != this.defaultTo) {
		delete this.defaultTo;
	}
}

/**
	registerLoading registers a loading.
	@private
	@this module:private/components/app/mail
	@param {!external:jQuery~Promise} promise - A promise representing the
	loading.
	@returns {Undefined}
*/
function registerLoading(promise) {
	this.progress.add(promise.catch(
		xhr => openError.call(this, alert.closable,
			client.error(xhr) || "どうしようもないエラーです")));
}

/**
	submit submits the mail.
	@param {!external:DOM~HTMLFormElement} target - A form representing
	properties of the mail.
	@returns {Undefined}
*/
function submit(target) {
	const param = {
		body:    target.body.value,
		subject: target.subject.value,
		to:      target.to.value,
	};

	const pickeds = this.picker.get();

	const iteration = pickeds.next();
	param.recipients = iteration.value.id;

	for (const picked of pickeds) {
		param.recipients = [param.recipients, picked.id].join(" ");
	}

	const inprogress = openInprogress.call(this, "送信しています…");

	this.progress.add(client.mail(param).then(response => {
		inprogress.remove();

		if (response.error == "mail_failure") {
			openError.call(this, alert.closable,
				"メールの送信に失敗しました");
		} else {
			openOK.call(this, alert.leavable, "送信しました");
		}
	}, xhr => {
		inprogress.remove();

		openError.call(this, alert.closable,
			client.error(xhr) || "どうしようもないエラーです");
	}));
}

export function oninit() {
	this.picker = picker.newPicker();
	this.picker.onloadstart = registerLoading.bind(this);
	this.picker.title = "誰に送る?";
	this.pickerLoading = this.picker.load();

	this.defaultTo = null;
	this.progress = new ProgressSum;

	registerLoading.call(this, this.pickerLoading);
	defaultTo.call(this);
}

export function onbeforeremove() {
	if (this.modalPicker) {
		this.modalPicker.remove();
	}
}

export function view() {
	const pickerBarrier = {
		pending:  "読み込み中です",
		rejected: "読み込みに失敗したため使用できません",
	}[this.pickerLoading.state()] || "";

	const submissionBarrier = pickerBarrier ||
		(this.picker.count() <= 0 ? "送る相手を選んでください" : "");

	return [
		m(progress, this.progress.html()),
		m(container, {style: {height: "100%"}},
			m("div", {className: "container", style: {height: "100%"}},
				m("form", {
					onsubmit: (function(event) {
						this(event.target);

						return false;
					}).bind(submit.bind(this)),

					style: {
						display:       "flex",
						flexDirection: "column",
						minHeight:     "100%",
					},
				},
					m("h1", "Mail"),
					m("div", {className: "form-group"},
						m("div", {
							className: "control-label",
							style:     {fontWeight: "bold"},
						}, "To"),
						m("div", {
							style: {
								display:  "flex",
								flexWrap: "wrap",
								width:    "100%",
							},
						},
							m("input", {
								className:   "form-control",
								name:        "to",
								onblur:      m.withAttr("value", updateTo.bind(this)),
								placeholder: "To",
								style:       {flex: "1"},
								value:       this.defaultTo,
							}), m("button", {
								className: "btn btn-default",
								disabled:  Boolean(pickerBarrier),
								onclick:   openPicker.bind(this),
								oncreate:  node => this.pickerLoading.done(() => {
									if (document.activeElement == document.body) {
										m.redraw(true);
										node.dom.focus();
									}
								}),
								title: pickerBarrier,
								type:  "button",
							},
								m("span", {"aria-hidden": "true"},
									m("span", {className: "glyphicon glyphicon-check"}),
									" "
								), "変更する"
							)
						)
					), m("label", {className: "form-group", style: {display: "block"}},
						m("div", {className: "control-label"}, "Subject"),
						m("input", {
							className:   "form-control",
							name:        "subject",
							placeholder: "Subject",
							style:       {fontWeight: "normal"},
						})
					), m("label", {
						className: "form-group",
						style:     {
							display:       "flex",
							flexDirection: "column",
							flex:          "1",
						},
					},
						m("div", {className: "control-label"}, "Body"),
						m("textarea", {
							className:   "form-control",
							name:        "body",
							placeholder: "Body",
							style:       {
								flex:       "1",
								fontWeight: "normal",
							},
						})
					), m("div", {className: "form-group"},
						m("button", {
							className: "btn btn-primary",
							disabled:  Boolean(submissionBarrier),
							title:     submissionBarrier,
						},
							m("span", {"aria-hidden": "true"},
								m("span", {className: "glyphicon glyphicon-send"}),
								" "
							), "送信"
						)
					)
				)
			)
		),
	];
}

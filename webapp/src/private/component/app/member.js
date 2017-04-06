/**
	@file member.js implements member component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/member */

/**
	module:private/component/app/member is a component to show a member.
	@name module:private/component/app/member
	@type !external:Mithril~Component
*/

import * as alert from "../alert";
import * as container from "../container";
import * as modal from "../../modal";
import * as primitive from "../member/primitive";
import * as progress from "../../progress";
import client from "../../client";

/**
	moveToPrimitiveConfirmation moves to the confirmation by the primitive.
	@private
	@this module:private/component/app/member
	@returns {Undefined}
*/
function moveToPrimitiveConfirmation() {
	this.primitive.confirm();
}

/**
	openXHRErrorAlert opens an error dialog. TODO
	@private
	@param TODO
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openXHRErrorAlert = (onclosed, error) => modal.add(
	alert.closable({onclosed},
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
			" "),
		client.error(error)));

/**
	openConfirmationOKAlert opens a dialog showing a successful result. TODO
	@private
	@param TODO
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openConfirmationOKAlert = onclosed => modal.add(
	alert.closable({onclosed},
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-ok"}),
			" "),
		"メールアドレスを確認しました。"));

/**
	openConfirmationError opens a dialog showing a error when confirming
	mail.
	@private
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openConfirmationError(onclosed) {
	let button;
	let modalNode;

	return modal.add({
		onmodalinit(node) {
			modalNode = node;
		},

		onmodalshown() {
			button.focus();
		},

		onmodalhidden: onclosed,

		view: () => m("div", {className: "modal-content"},
			m("div", {
				className: "modal-body",
				id:        "component-app-member-modal",
			}, "メールアドレスの確認に失敗しました。たぶん時間切れとかそんなところじゃないですかね? 確認メールを再送信しますか?"),
			m("div", {className: "modal-footer"},
				m("button", {
					"data-dismiss": "modal",
					className:      "btn btn-default",
				}, "やっぱやめる"),
				m("button", {
					className: "btn btn-primary",

					onclick: () => {
						moveToPrimitiveConfirmation.call(this);
						modalNode.remove();
					},

					oncreate(node) {
						button = node.dom;
					},
				}, "再送信する"))),
	});
}

/**
	openConfirmationBusyAlert opens a dialog showing progress. TODO
	@private
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openConfirmationBusyAlert =
	() => modal.add({backdrop: "static"}, alert.busy("送信しています…"));

export function oninit() {
	if (m.route.param("id") == client.getID()) {
		const confirm = m.route.param("confirm");
		if (confirm) {
			const confirmationModal = openConfirmationBusyAlert.call(this);
			const confirmationProgress = progress.add({
				"aria-describedby": alert.bodyID,
				value:              0,
			});

			client.confirm(confirm).then(() => {
				confirmationModal.remove();
				openConfirmationOKAlert.call(this, confirmationProgress.remove);
			}, error => {
				confirmationModal.remove();

				if (error == "invalid_request") {
					confirmationProgress.updateARIA(
						{"aria-describedby": "component-app-member-modal"});

					openConfirmationError.call(this, confirmationProgress.remove);
				} else {
					openXHRErrorAlert.call(this, confirmationProgress.remove, error);
				}
			}, event => confirmationProgress.updateValue(
				{max: event.total, value: event.loaded}));

			return;
		}
	}

	this.primitive = new primitive.State;
}

export function oncreate() {
	this.primitive.start(m.route.param("id"));
	this.primitive.focus();
}

export function onbeforeremove() {
	this.primitive.end();
}

export function view() {
	return this.primitive && m(container,
		m("div", {
			className: "container",
			style:     {textAlign: "center"},
		},
			m(this.primitive.getForm() ? "form" : "div",
				m("div",
					m("div", {
						"aria-hidden": (this.error == null).toString(),
						style:         {
							display:   "inline-block",
							minHeight: "8rem",
							textAlign: "left",
						},
					},
						this.error && m("div", {
							className: "alert alert-danger",
							role:      "alert",
						}, m(primitive.error, this.error)),
						m("h1", {style: {fontSize: "x-large"}},
							m(this.primitive.header)),
						m(this.primitive.body))),
				m(this.primitive.button))));
}

/**
	@file member.js implements member component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/member */

/**
	module:private/components/app/member is a component to show a member.
	@name module:private/components/app/member
	@type !external:Mithril~Component
*/

import * as alert from "../alert";
import * as container from "../container";
import * as modal from "../../modal";
import * as primitive from "../member/primitive";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";
import client from "../../client";

/**
	registerLoading registers a promise describing a loading.
	@private
	@this module:private/components/app/member
	@param {!external:jQuery~Promise} promise - The promise describing a
	loading.
	@returns {Undefined}
*/
function registerLoading(promise) {
	this.progress.add(promise.then(submitted => {
		if (submitted) {
			delete this.error;
		}
	}, message => {
		this.error = message;
	}));
}

/**
	load loads the remote content.
	@private
	@this module:private/components/app/member
	@returns {Undefined}
*/
function load() {
	this.loaded = false;
	this.primitive = new primitive.State;

	this.primitive.setOnload(() => {
		this.loaded = true;

		if (modal.isEmpty()) {
			this.primitive.focus();
		}
	});

	this.primitive.setOnloadstart(registerLoading.bind(this));
	this.primitive.setID(m.route.param("id"));
}

/**
	moveToPrimitiveConfirmation moves to the confirmation by the primitive.
	@private
	@this module:private/components/app/member
	@returns {Undefined}
*/
function moveToPrimitiveConfirmation() {
	this.primitive.confirm();
}

/**
	openError opens an error dialog.
	@private
	@param {...?external:Mithril~Children} children - An error message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openError() {
	return modal.unshift(alert.closable(
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
			" "
		), ...arguments
	));
}

/**
	openOK opens a dialog showing a successful result.
	@private
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openOK() {
	return modal.unshift(alert.closable(
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-ok"}),
			" "
		), ...arguments
	));
}

/**
	openConfirmationError opens a dialog showing a error when confirming
	mail.
	@private
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openConfirmationError() {
	let button;
	let modalNode;

	return modal.unshift({
		onmodalinit(node) {
			modalNode = node;
		},

		onmodalshown() {
			button.focus();
		},

		view: () => m("div", {className: "modal-content"},
			m("div", {className: "modal-body"},
				"メールアドレスの確認に失敗しました。たぶん時間切れとかそんなところじゃないですかね? 確認メールを再送信しますか?"
			), m("div", {className: "modal-footer"},
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
				}, "再送信する")
			)
		),
	});
}

/**
	openInprogress opens a dialog showing progress.
	@private
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openInprogress() {
	return modal.unshift({backdrop: "static"},
		alert.inprogress(...arguments));
}

export function oninit() {
	this.progress = new ProgressSum;

	if (m.route.param("id") == client.getID()) {
		const confirm = m.route.param("confirm");
		if (confirm) {
			const inprogress = openInprogress.call(this, "送信しています…");

			this.progress.add(client.userConfirm(confirm).then(() => {
				inprogress.remove();
				openOK.call(this, "メールアドレスを確認しました。");
				load.call(this);
			}, xhr => {
				inprogress.remove();

				if (xhr.responseJSON && xhr.responseJSON.error == "invalid_request") {
					openConfirmationError.call(this);
				} else {
					openError.call(this,
						client.error(xhr) || "どうしようもないエラーです。");
				}

				load.call(this);
			}));

			return;
		}
	}

	load.call(this);
}

export function view() {
	return [
		m(progress, this.progress.html()),
		this.primitive && m(container,
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
								m(this.primitive.header)
							), m(this.primitive.body)
						)
					),
					m(this.primitive.button)
				)
			)
		),
	];
}

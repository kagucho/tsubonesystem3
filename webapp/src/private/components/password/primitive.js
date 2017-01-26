/**
	@file primitive.js provides the primitive elements of the UI to update
	the password of the user.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/password/primitive */

import * as client from "../../client";

class FormControl {
	constructor(label) {
		this.label = label;
	}

	updateValidity(target) {
		const formControl = this.formControls[target.name];

		if (!formControl.validity) {
			this.invalids++;
		}

		formControl.validity = target.validationMessage;
	}

	checkValidity(target) {
		const formControl = this.formControls[target.name];

		if (target.checkValidity() && formControl.validity) {
			delete formControl.validity;
			this.invalids--;
		}
	}

	dismissValidity(target) {
		const formControl = this.formControls[target.name];

		if (formControl.validity) {
			formControl.checkValidity.call(this, target);
		}
	}
}

const handleEventTarget = handler => (function(event) {
	this(event.target);
}).bind(handler);

export class controller {
	constructor(onloadstart) {
		this.fireLoadstart = onloadstart || $.noop;

		this.formControls = {
			oldPassword: new class extends FormControl {
				constructor() {
					super("現在のパスワードを入力してください");
				}
			},

			newPassword: new class extends FormControl {
				constructor() {
					super("新しいパスワードを入力してください");
				}

				checkValidity(target) {
					if (this.formControls.verificationPassword.validity) {
						this.formControls.verificationPassword.checkValidity.call(this,
							target.form.verificationPassword);
					}

					super.checkValidity(target);
				}

				dismissValidity(target) {
					if (target.form.elements.verificationPassword.validity.customError) {
						this.formControls.verificationPassword.checkValidity.call(this,
							target.form.elements.verificationPassword);
					}

					super.dismissValidity(target);
				}
			},

			verificationPassword: new class extends FormControl {
				constructor() {
					super("新しいパスワードをもう一度入力してください");
				}

				checkValidity(target) {
					if (target.value != target.form.newPassword.value) {
						target.setCustomValidity("違うよ");

						return super.updateValidity(target);
					}

					target.setCustomValidity("");
					super.checkValidity(target);
				}
			},
		};

		this.invalids = 0;
	}

	submit(target) {
		this.formControls.verificationPassword.checkValidity.call(this,
			target.form.elements.verificationPassword);

		if (!target.form.checkValidity()) {
			return;
		}

		this.fireLoadstart(client.memberUpdatePassword({
			"old": target.form.elements.oldPassword.value,
			"new": target.form.elements.newPassword.value,
		}).catch(xhr => {
			throw client.error(xhr) || "どうしようもないエラーが発生しました。";
		}));
	}
}

export const titleView = "パスワード変更";

export const bodyView = control => [
	"パスワードを変更します。",
	m("div", $.map(control.formControls, (value, key) => m("label",
		{style: {display: "block", margin: "1rem"}},
		m("div", {className: "control-label"}, value.label),
		m("div", {className: "component-password-record"},
			m("input", {
				className: "form-control",
				maxlength: "28",
				name:      key,
				onblur:    handleEventTarget(value.checkValidity.bind(control)),
				oninput:   handleEventTarget(value.dismissValidity.bind(control)),
				oninvalid: handleEventTarget(value.updateValidity.bind(control)),
				required:  true,
				style:     {
					display:  "inline-block",
					margin:   "1rem",
					maxWidth: "40ch",
					position: "static",
				},
				type: "password",
			}), m("div", {
				style: {
					display:  "inline-block",
					minWidth: "32ch",
				},
			}, value.validity && m("div", {
				className: "alert alert-danger",
				role:      "alert",
				style:     {
					display: "inline-block",
					padding: "0.5rem",
					margin:  "0.5rem",
				},
			},
				m("span", {ariaHidden: "true"},
					m("span", {
						className: "glyphicon glyphicon-exclamation-sign",
					}), " "
				), value.validity
			))
		)
	))),
];

export function buttonView(control) {
	const attributes = {
		className: "btn btn-primary",
		onclick:   handleEventTarget(control.submit.bind(control)),
		type:      "button",
	};

	if (control.invalids) {
		attributes.disabled = true;
		attributes.title = "不正な項目があります。";
	} else if (control.progress && control.progress.max != control.progress.value) {
		attributes.disabled = true;
		attributes.title = "まだ準備できていません。";
	} else {
		attributes.disabled = false;
		attributes.title = "";
	}

	return m("button", attributes, "送信");
}

/**
	errorView returns the view of a given error message.
	@param {?external:ES.String} message - The error message.
	@returns {!external:Mithril~Children} The view.
*/
export function errorView(message) {
	return [
		m("span", {ariaHidden: "true"},
			m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
			" "
		), message,
	];
}

export const inprogressView = [
	m("span", {ariaHidden: "true"},
		m("span", {className: "glyphicon glyphicon-hourglass"}),
		" "
	), "送信しています…",
];

export const successView = [
	m("span", {ariaHidden: "true"},
		m("span", {className: "glyphicon glyphicon-ok"}),
		" "
	), "送信しました",
];

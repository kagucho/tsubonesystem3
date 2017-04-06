/**
	@file primitive.js provides the primitive elements of the UI to update
	the password of the user.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/password/primitive */

import * as alert from "../alert";
import * as modal from "../../modal";
import * as progress from "../../progress";
import client from "../../client";

/**
	Input is a class to hold a state of an input element.
	@private
	@extends Object
*/
class Input {
	/**
		constructor constructs module:private/component/password/primitive~Input
		@param {!String} autocomplete - The value to be set to
		autocomplete attribute.
		@param {!String} label - The string of the label.
		@returns {Undefined}
	*/
	constructor(autocomplete, label) {
		this.autocomplete = autocomplete;
		this.label = label;
	}

	/**
		bindToEvent returns a function to handle event.
		@function
		@param {!module:private/component/password/primitive~Internal}
		internal - The internal state to be bound.
		@param {!String} key - The key for the function to be bound to.
		@returns {!module:private/component/password/primitive~HandleEventTarget}
		A function to handle event.
	*/
	static bindToEvent(internal, key) {
		return function(event) {
			this(event.target);
		}.bind(this[key].bind(internal));
	}

	/**
		updateValidity updates the validity. Note the value of this.
		@this module:private/component/password/primitive~Internal
		@param {!external:DOM~HTMLInputElement} target - An
		HTMLInputElement representing the validity.
		@returns {Undefined}
	*/
	static updateValidity(target) {
		const input = this.inputs[target.name];

		if (!input.validity) {
			this.invalids++;
		}

		input.validity = target.validationMessage;
	}

	/**
		checkValidity checks the validity of the value represented by
		the given element and reflects the result to the element.
		Note the value of this.
		@this module:private/component/password/primitive~Internal
		@param {!external:DOM~HTMLInputElement} target - An
		HTMLInputElement representing the value and the validity.
		@returns {Undefined}
	*/
	static checkValidity(target) {
		const input = this.inputs[target.name];

		if (target.checkValidity() && input.validity) {
			delete input.validity;
			this.invalids--;
		}
	}

	/**
		reportValidity checks the validity of the value represented by
		the given element and reports the result to the user.
		Note the value of this.
		@this module:private/component/password/primitive~Internal
		@param {!external:DOM~HTMLInputElement} target - An
		HTMLInputElement representing the value and the validity.
		@returns {Undefined}
	*/
	static reportValidity(target) {
		const input = this.inputs[target.name];

		if (target.reportValidity() && input.validity) {
			delete input.validity;
			this.invalids--;
		}
	}

	/**
		dismissValidity dismisses the validity and introduces the next
		validity if exists. Note the value of this.
		@this module:private/component/password/primitive~Internal
		@param {!external:DOM~HTMLInputElement} target - An
		HTMLInputElement representing the value and the validity.
		@returns {Undefined}
	*/
	static dismissValidity(target) {
		const input = this.inputs[target.name];

		if (input.validity) {
			input.constructor.checkValidity.call(this, target);
		}
	}

	/**
		autocomplete is the string to be set to autocomplete attribute.
		@member {!String} module:private/component/password/primitive~Input#autocomplete
	*/

	/**
		label is the string of the label.
		@member {!String} module:private/component/password/primitive~Input#label
	*/

	/**
		validity is the string representing the validity. If it is null,
		it is valid. Otherwise, it is invalid and validity represents an
		advice to make it valid.
		@member {?String} module:private/component/password/primitive~Input#validity
	*/
}

/**
	Internal is a class to hold the internal state.
	@private
	@extends Object
*/
class Internal {
	/**
		constructor constructs module:private/component/password/primitive~Internal.
		@returns {Undefined}
	*/
	constructor(onemptied) {
		this.inputs = {
			current: new class extends Input {
				constructor() {
					super("current-password",
						"現在のパスワードを入力してください");
				}
			},

			new: new class extends Input {
				constructor() {
					super("new-password",
						"新しいパスワードを入力してください");
				}

				static dismissValidity(target) {
					this.inputs.verification.constructor.dismissValidity.call(this,
						target.form.verification);

					super.dismissValidity(target);
				}
			},

			verification: new class extends Input {
				constructor() {
					super("new-password",
						"新しいパスワードをもう一度入力してください");
				}

				static checkValidity(target) {
					if (target.value != target.form.new.value) {
						target.setCustomValidity("違うよ");

						return super.updateValidity(target);
					}

					target.setCustomValidity("");
					super.checkValidity(target);
				}

				static reportValidity(target) {
					target.setCustomValidity(target.value == target.form.new.value ?
						"" : "違うよ");

					super.reportValidity(target);
				}
			},
		};

		this.invalids = 0;
		this.onemptied = onemptied;
	}

	/**
		prepareFocus prepares for focusing.
		@param {!external:Mithril~Node} node - The node of the element
		to be focused.
		@returns {Undefined}
	*/
	prepareFocus(node) {
		this.focus = node.dom;
	}

	/**
		submit submits a request to update the password of the user.
		@param {!external:DOM~HTMLFormElement} target - A form
		representing the values and the validities.
		@returns {Undefined}
	*/
	submit(target) {
		this.inputs.verification.constructor.reportValidity.call(this, target.verification);

		if (!target.reportValidity()) {
			return;
		}

		const submissionModal = modal.add(alert.busy("送信しています…"));

		const submissionProgress = progress.add({
			"aria-describedby": alert.bodyID,
			value:              0,
		}, true);

		client.patchUser({
			current_password: target.current.value,
			new_password:     target.new.value,
		}).then(() => {
			submissionModal.remove();

			modal.add(this.onemptied ?
				alert.closable({
					onclosed() {
						submissionProgress.remove();
						this.onemptied();
					},
				}, "送信しました。") :
				alert.leavable({
					onclosed: submissionProgress.remove,
				}, "送信しました。"));
		}, error => {
			submissionModal.remove();

			modal.add(alert.closable(
				{onclosed: submissionProgress.remove},
				error == "invalid_request" ?
					"パスワードが違います" :
					client.error(error)));
		}, event => submissionProgress.updateValue(
			{max: event.total, value: event.loaded}));
	}

	/**
		bodyView returns a view of the body.
		@returns {!external:Mithril~Children} A view of the body.
	*/
	bodyView() {
		const {inputs} = this;

		return [
			"パスワードを変更します。",
			m("div", ["current", "new", "verification"].map(key => m("label",
				{style: {display: "block", margin: "1rem"}},
				m("div", {className: "control-label"}, inputs[key].label),
				m("div", {className: "component-password-record"},
					m("input", {
						autocomplete: inputs[key].autocomplete,
						className:    "form-control",
						inputmode:    "verbatim",
						maxlength:    "128",
						name:         key,
						oncreate:     key == "current" &&
							this.prepareFocus.bind(this),
						onchange:  inputs[key].constructor.bindToEvent(this, "checkValidity"),
						oninput:   inputs[key].constructor.bindToEvent(this, "dismissValidity"),
						oninvalid: inputs[key].constructor.bindToEvent(this, "updateValidity"),
						pattern:   "[ -~]*",
						required:  true,
						style:     {
							display:  "inline-block",
							margin:   "1rem",
							maxWidth: "32ch",
							position: "static",
						},
						title: "ASCII (英数字と空白, 様々な記号を含む. 詳しくはググれ.)",
						type:  "password",
					}), m("div", {style: {display: "inline-block"}},
						inputs[key].validity && m("div", {
							className: "alert alert-danger",
							role:      "alert",
							style:     {
								display: "inline-block",
								padding: "0.5rem",
								margin:  "0.5rem",
							},
						},
							m("span", {"aria-hidden": "true"},
								m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
								" "
							), inputs[key].validity
						)
					)
				)
			))),
		];
	}

	/**
		buttonView returns a view of the button.
		@returns {!external:Mithril~Children} A view of the button.
	*/
	buttonView() {
		const attrs = {
			className: "btn btn-primary",
			onclick:   event => {
				this.submit(event.target.form);

				return false;
			},
		};

		if (this.invalids) {
			attrs.disabled = true;
			attrs.title = "そんな入力内容で大丈夫か?";
		} else if (this.progress && this.progress.max != this.progress.value) {
			attrs.disabled = true;
			attrs.title = "もう少し待つのだぞ。";
		}

		return m("button", attrs, "送信");
	}

	/**
		inputs is a list of module:private/component/password/primitive~Input.
		@member {!Object.<!String, !module:private/component/password/primitive~Input>}
		module:private/component/password/primitive~Internal#inputs
	*/

	/**
		invalids is the number of invalid inputs.
		@member {!Number} module:private/component/password/primitive~Internal#invalids
	*/
}

/**
	newState returns a new State.
	@param TODO
	@returns {!module:private/component/password/primitive~State} A new
	State.
*/
export function newState(onemptied) {
	const internal = new Internal(onemptied);

	return {
		body: {
			view() {
				return internal.bodyView();
			},
		},

		button: {
			view() {
				return internal.buttonView();
			},
		},

		focus() {
			internal.focus.focus();
		},
	};
}

/**
	title is the string of the title.
	@type !String
*/
export const title = "パスワード変更";

/**
	HandleEventTarget is a function to handle event.
	@private
	@callback module:private/component/password/primitive~HandleEventTarget
	@param {!external:DOM~Event} event - An event to handle.
	@returns {Undefined}
*/

/**
	State is an interface for the state.
	@interface module:private/component/password/primitive~State
	@extends Object
*/
/**
	body is a component to draw the body.
	@member {!external:Mithril~Component}
		module:private/component/password/primitive~State#body
*/

/**
	button is a component to draw the button.
	@member {!external:Mithril~Component}
		module:private/component/password/primitive~State#button
*/

/**
	focus focuses.
	@function module:private/component/password/primitive~State#focus
	@returns {Undefined}
*/

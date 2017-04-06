/**
	@file modal.js implements a component to show modal dialogs to update
	the password of the user.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/password/modal */

import * as primitive from "./primitive";

/**
	labelID is the ID of the labelling element.
	@type !String
*/
export const labelID = "component-password-modal-title";

/**
	newComponent returns a new component to show a modal dialog to update
	the password of the user.
	@returns {!module:private/modal~Component} A new component to show a
	modal dialog to update the password of the user.
*/
export function newComponent() {
	const state = primitive.newState();

	return {
		onmodalshown() {
			state.focus();
		},

		view() {
			return m("form", {className: "modal-content"},
				m("div", {className: "modal-header"},
					m("button", {
						"aria-label":   "閉じる",
						"data-dismiss": "modal",
						className:      "close",
						type:           "button",
					}, m("span", {"aria-hidden": "true"}, "×")),
					m("a", {
						className: "lead modal-title",
						href:      "#!password",
						id:        labelID,
					}, primitive.title)
				), m("div", {className: "modal-body"},
					m(state.body)
				), m("div", {className: "modal-footer"},
					m(state.button)
				)
			);
		},
	};
}

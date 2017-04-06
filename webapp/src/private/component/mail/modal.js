/**
	@file modal.js implements a component to show modal dialogs to show or
	create email.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/mail/modal */

import * as primitive from "./primitive";

/**
	TODO
*/
export const labelID = "component-mail-modal-title";

/**
	TODO
*/
export default subject => {
	const state = new primitive.default(subject);
	let button;

	return {
		onmodalinit(node) {
			state.onemptied = node.remove;
		},

		onmodalshow() {
			state.start();
		},

		onmodalshown() {
			button.focus();
		},

		onmodalhide() {
			state.end();
		},

		view() {
			return m(subject ? "div" : "form",
				{className: "modal-content"},
				m("div", {className: "modal-header"},
					m("button", {
						"aria-label":   "閉じる",
						"data-dismiss": "modal",
						className:      "close",

						oncreate(node) {
							button = node.dom;
						},

						type: "button",
					}, m("span", {ariaHidden: "true"}, "×")),
					m("a", {
						className: "lead modal-title",
						href:      "#!mail?subject=" + subject,
						id:        labelID,
					}, primitive.title),
					m("div", {
						className: "modal-body",
						style:     {maxWidth: "80ch"},
					}, m(state.body)),
					m("div", {
						className: "modal-footer",
						style:     {
							display:        "flex",
							justifyContent: "space-around",
						},
					},
						m("button", {
							"data-dismiss": "modal",
							className:      "btn btn-danger",
						}, "閉じる"),
						m(state.button))));
		},
	};
};

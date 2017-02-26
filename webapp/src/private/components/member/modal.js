/**
	@file modal.js implements a component to show modal dialogs describing
	members.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/member/modal */

import * as primitive from "./primitive";

/**
	module:private/components/member/modal returns a component to draw
	a modal dialog showing a member.
	@param {?String} id - The ID. If it is null, it will be the form to
	create a new member. Otherwise, it will describe the member identified
	with the ID.
	@param {?module:private/components/member/primitive~Onloadstart}
	onloadstart - The function to be called after a loading starts.
	@returns {!module:private/modal~Component} A component to draw a modal
	dialog showing a member.
*/
export default (id, onloadstart) => {
	const state = new primitive.State;
	let button;

	state.setOnloadstart(onloadstart);

	return {
		onmodalinit(node) {
			state.setOnemptied(node.remove);
			state.setID(id);
		},

		onmodalshown() {
			button.focus();
		},

		view() {
			return m(state.getForm() ? "form" : "div",
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
						href:      "#!member?id=" + id,
						id:        "component-member-modal-title",
					}, m(state.header)),
					m("div", {
						className: "modal-body",
						style:     {
							maxWidth: "80ch",
						},
					}, m(state.body)),
					m("div", {className: "modal-footer"},
						m(state.button))
				)
			);
		},
	};
};

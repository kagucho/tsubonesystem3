/**
	@file modal.js implements a component to show modal dialogs to update
	the password of the user.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/password/modal */

import * as alert from "../alert";
import * as modal from "../../modal";
import * as primitive from "./primitive";

/**
	registerLoading registers a loading.
	@private
	@param {!module:private/modal~Node} node - The node in the list of the
	modal dialogs.
	@param {!external:jQuery~Promise} promise - A promise representing a
	loading.
	@returns {Undefined}
*/
function reflectLoadingToModal(node, promise) {
	const loading = modal.unshift(alert.inprogress, primitive.inprogress);

	const filteredPromise = promise.then(submission => {
		loading.remove();
		node.remove();

		modal.unshift((submission ? alert.leavable : alert.closable)(
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-ok"}),
				" "
			), primitive.success
		));
	}, xhr => {
		loading.remove();

		modal.unshift(alert.closable(
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
				" "
			), primitive.error(xhr)
		));
	});

	return filteredPromise;
}

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
	let modalNode;

	return {
		oninit() {
			state.setOnloadstart(
				reflectLoadingToModal.bind(modalNode));
		},

		onmodalinit(initModalNode) {
			modalNode = initModalNode;
		},

		onshown() {
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

/**
	@file alert.js implements alert component for modal dialogs.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/alert */

import * as navigator from "../navigator";

/**
	TODO
*/
export const bodyID = "component-alert-body";

/**
	busy returns a component to draw an alert to indicate a work is in
	busy with the given message.
	@param {...?external:Mithril~Children} children - The message indicating
	the progress.
	@returns {!module:private/modal~Component} A component to draw a modal
	dialog to indicate progress with the given message.
*/
export function busy() {
	return {
		view: () => m("div", {className: "modal-content"},
			m("div", {className: "modal-body", id: bodyID},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-hourglass"}),
					" "),
				...arguments)),
	};
}

/**
	closable returns a component to draw a closable alert.
	@param {!module:private/component/alert~ClosableAttrs} [attrs] - A set
	of attributes for the component.
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Component} A component to draw a
	closable alert.
*/
export function closable() {
	let button;
	let childrenIndex;
	let onclosed;

	if (arguments[0] && typeof arguments[0].onclosed == "function") {
		childrenIndex = 1;
		([{onclosed}] = arguments);
	} else {
		childrenIndex = 0;
		onclosed = null;
	}

	return {
		onmodalshown() {
			button.focus();
		},

		onmodalremove: onclosed,

		view: () => m("div", {className: "modal-content"},
			m("div", {className: "modal-body", id: bodyID},
				Array.prototype.slice.call(arguments, childrenIndex)),
			m("div", {className: "modal-footer"},
				m("button", {
					"data-dismiss": "modal",
					className:      "btn btn-default",

					oncreate: node => {
						button = node.dom;
					},
				}, "閉じる"))),
	};
}

/**
	leavable returns a component to draw an alert which will leave the
	current page after it gets closed.
	@param {!module:private/component/alert~ClosableAttrs} [attrs] - A set
	of attributes for the component.
	@param {...?external:Mithril~Children} children - The message indicating
	the progress.
	@returns {!module:private/modal~Component} A component to draw a
	leavable alert.
*/
export function leavable() {
	let button;
	let leave;
	let childrenIndex;
	let onclosed;

	if (arguments[0] && typeof arguments[0].onclosed == "function") {
		childrenIndex = 1;
		([{onclosed}] = arguments);
	} else {
		childrenIndex = 0;
		onclosed = $.noop;
	}

	return {
		onmodalshown() {
			button.focus();
		},

		onmodalremove() {
			onclosed();
			leave();
		},

		oninit() {
			leave = navigator.leaver();
		},

		view: () => {
			let hideView;
			switch (leave) {
			case navigator.top:
				hideView = m("a", {
					"data-dismiss": "modal",
					className:      "btn btn-default",
					href:           "",
					oncreate:       node => button = node.dom,
					role:           "button",
				}, "トップページへ");
				break;

			case navigator.back:
				hideView = m("button", {
					"data-dismiss": "modal",
					className:      "btn btn-default",
					oncreate:       node => button = node.dom,
				}, "戻る");
				break;

			default:
				throw new Error("unknown leaver");
			}

			return m("div", {className: "modal-content"},
				m("div", {className: "modal-body", id: bodyID},
					Array.prototype.slice.call(arguments, childrenIndex)),
				m("div", {className: "modal-footer"}, hideView));
		},
	};
}

/**
	ClosableAttrs is a set of attributes for closable.
	@typedef module:private/component/alert~ClosableAttrs
	@property {?module:private/modal~LifecycleMethod} onclosed - A hook to
	be called after the closable gets closed.
*/

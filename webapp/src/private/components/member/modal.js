/**
	@file modal.js implements a component to show modal dialogs describing
	members.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/member/modal */

/**
	module:private/components/member/modal is a component to show modal
	dialogs describing members.
	This component is reusable but NOT reentrant.
	@name module:private/components/member/modal
	@type external:Mithril~Component
*/

import * as primitive from "./primitive";
import {ensureRedraw} from "../../mithril";

export class controller {
	constructor() {
		this.callbacks = {};
		this.memberAttributes = {
			onemptied: () => {
				this.finished = true;
			},

			onloadstart: promise => this.callbacks.onloadstart(
				promise.catch(this.errors.push.bind(this.errors))),

			onmodalshow: () => {
				this.shown.next = this.shown.current;
				this.shown.current = "modal";
			},

			onmodalhide: this.restoreShown.bind(this),
		};

		this.primitive = new primitive.controller(this.memberAttributes);
		this.errors = [];
		this.shown = {};
	}

	finishBody() {
		if (this.shown.current == "body") {
			this.callbacks.onhidden();
		}
	}

	restoreShown() {
		this.shown.current = this.shown.next;
		this.shown.next = "body";
		if (this.finished) {
			this.finishBody();
		}
	}

	shiftError() {
		this.errors.shift();
		if (!this.errors.length) {
			this.restoreShown();
		}
	}

	updateAttributes(attributes) {
		for (const key of ["onhidden", "onloadstart"]) {
			this.callbacks[key] = attributes[key] || $.noop;
		}

		this.memberAttributes.id = attributes.id;
		this.primitive.updateAttributes(this.memberAttributes);

		if (!this.shown.current) {
			this.shown.current = "body";
			this.shown.next = "body";
		}

		if (this.errors.length) {
			this.shown.current = "error";
		}
	}
}

export function view(control, attributes) {
	control.updateAttributes(attributes);

	const body = [
		m("div", {
			ariaHidden:     (control.shown.current != "body").toString(),
			ariaLabelledby: "component-member-modal-title",
			className:      "modal fade",
			config:         (function(element, initialized, context) {
				const jquery = $(element);

				if (!initialized) {
					context.onunload = jquery.modal.bind(jquery, "hide");

					jquery.on("hidden.bs.modal",
						ensureRedraw.bind(undefined, this.finishBody.bind(this)));
				}

				if (this.primitive.critical) {
					jquery.modal({show: false, backdrop: "static"});
				}

				jquery.modal(this.shown.current == "body" ? "show" : "hide");
			}).bind(control),
			role:     "dialog",
			tabindex: "-1",
		}, m("div", {className: "modal-dialog", role: "document"},
			m("form", {className: "modal-content"},
				m("div", {className: "modal-header"},
					m("button", {
						ariaLabel:      "閉じる",
						className:      "close",
						type:           "button",
						"data-dismiss": "modal",
					}, m("span", {ariaHidden: "true"}, "×")),
					m("a", {
						className: "lead modal-title",
						href:      "#!member?id=" + control.memberAttributes.id,
						id:        "component-member-modal-title",
					}, primitive.headerView(control.primitive))
				), m("div", {className: "modal-body"},
					primitive.bodyView(control.primitive)
				), m("div", {className: "modal-footer"},
					primitive.buttonView(control.primitive)
				)
			)
		)), primitive.modalView(control.primitive),
	];

	if (control.errors.length) {
		body.push(m("div", {
			ariaHidden: (control.shown.current != "error").toString(),
			className:  "modal fade",
			config:     (function(element, initialized, context) {
				const jquery = $(element);

				if (!initialized) {
					context.onunload = jquery.modal.bind(jquery, "hide");

					jquery.on("hidden.bs.modal",
						() => this.shown.current == "error" &&
							ensureRedraw(this.shiftError.bind(this)));
				}

				jquery.modal("show");
			}).bind(control),
			role:     "dialog",
			tabindex: "-1",
		}, m("div", {className: "modal-dialog", role: "document"},
			m("div", {className: "modal-content"},
				m("div", {className: "modal-body"},
					primitive.errorView(control.errors[0])),
				m("div", {className: "modal-footer"},
					m("button", {
						className:      "btn btn-default",
						type:           "button",
						"data-dismiss": "modal",
					}, "閉じる")
				)
			)
		)));
	}

	return m("div", body);
}

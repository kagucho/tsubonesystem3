/**
	@file modal.js implements a component to show modal dialogs describing
	members.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module private/components/member/modal */

/**
	module:private/components/member/modal is a component to show modal
	dialogs describing members.
	@name module:private/components/member/modal
	@type external:Mithril~Component
*/

import * as primitive from "./primitive";
import {ensureRedraw} from "../../mithril";

export class controller {
	constructor() {
		this.memberAttributes = {
			onerror:   message => this.messages.push({body: message, type: "error"}),
			onsuccess: message => this.messages.push({body: message, type: "success"}),

			onemptied: () => {
				this.finished = true;
			},

			onloadend: submission => {
				if (submission) {
					this.finished = true;
				}

				this.onloadendCallback();
			},

			onmodalshow: () => {
				this.shown.next = this.shown.current;
				this.shown.current = "modal";
			},

			onmodalhide: this.restoreShown.bind(this),
		};

		this.primitive = new primitive.controller(this.memberAttributes);
		this.messages = [];
		this.shown = {};
	}

	onhidden() {
		if (this.shown.current == "body") {
			this.onhiddenCallback();
		}
	}

	restoreShown() {
		this.shown.current = this.shown.next;
		this.shown.next = "body";
		if (this.finished) {
			this.onhidden();
		}
	}

	shiftMessage() {
		this.messages.shift();
		if (!this.messages.length) {
			this.restoreShown();
		}
	}

	updateAttributes(attributes) {
		for (const key of ["id", "onloadstart", "onprogress"]) {
			this.memberAttributes[key] = attributes[key];
		}

		this.onhiddenCallback = attributes.onhidden;
		this.onloadendCallback = attributes.onloadend;

		this.primitive.updateAttributes(this.memberAttributes);

		if (!this.shown.current) {
			this.shown.current = "body";
			this.shown.next = "body";
		}

		if (this.messages.length) {
			this.shown.current = "message";
		}
	}
}

export function view(control, attributes) {
	control.updateAttributes(attributes);

	const body = [
		m("div", {
			ariaHidden:     (control.shown.current != "body").toString(),
			ariaLabelledby: "member-modal-title",
			className:      "modal fade",
			config:         (function(element, initialized, context) {
				if (!initialized) {
					context.onunload = (function() {
						$(this).modal("hide");
					}).bind(element);

					$(element).on("hidden.bs.modal",
						this.onhidden.bind(this));
				}

				if (this.primitive.critical) {
					$(element).modal({show: false, backdrop: "static"});
				}

				$(element).modal(this.shown.current == "body" ? "show" : "hide");
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
						id:        "member-modal-title",
					}, primitive.headerView(control.primitive))
				), m("div", {className: "modal-body"},
					primitive.bodyView(control.primitive)
				), m("div", {className: "modal-footer"},
					primitive.buttonView(control.primitive)
				)
			)
		)), primitive.modalView(control.primitive),
	];

	if (control.messages.length) {
		body.push(m("div", {
			ariaHidden:     (control.shown.current != "message").toString(),
			ariaLabelledBy: "member-modal-message-title",
			className:      "modal fade",
			config:         (function(element, initialized, context) {
				if (!initialized) {
					context.onunload = (function() {
						$(this).modal("hide");
					}).bind(element);

					$(element).on("hidden.bs.modal", () => {
						if (this.shown.current == "message") {
							ensureRedraw(this.shiftMessage.bind(this));
						}
					});

					$(element).modal("show");
				}
			}).bind(control),
			role:     "dialog",
			tabindex: "-1",
		}, m("div", {className: "modal-dialog", role: "document"},
			m("div", {className: "modal-content"},
				m("div", {className: "modal-body"},
					{
						error:   primitive.errorView,
						success: primitive.successView,
					}[control.messages[0].type](control.messages[0].body)
				), m("div", {className: "modal-footer"},
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

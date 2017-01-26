/**
	@file modal.js implements a component to show modal dialogs to update
	the password of the user.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/password/modal */

/**
	module:private/components/password/modal is a component to show modal
	dialogs to update the password of the user.
	This component is reusable but NOT reentrant.
	@name module:private/components/password/modal
	@type external:Mithril~Component
*/

import * as primitive from "./primitive";
import {ensureRedraw} from "../../mithril";

export class controller {
	constructor() {
		this.primitive = new primitive.controller(promise => {
			this.fireLoadstart(promise);
			this.loading = true;

			promise.then(() => {
				this.slides.push({type: "success"});
				delete this.loading;
			}, message => {
				this.slides.push({type: "error", message});
				delete this.loading;
			});
		});

		this.slides = [];
	}

	finishSlide() {
		const {type} = this.slides.shift();

		if (type == "success") {
			this.finish();
		}
	}

	updateAttributes(attributes) {
		this.finish = attributes.onhidden || $.noop;
		this.fireLoadstart = attributes.onloadstart;
	}
}

export function view(control, attributes) {
	let content;

	control.updateAttributes(attributes);

	if (control.slides.length) {
		switch (control.slides[0].type) {
		case "error":
			content = m("div", {className: "modal-content"},
				m("div", {className: "modal-body"},
					primitive.errorView(control.slides[0].message)
				), m("div", {className: "modal-footer"},
					m("button", {
						className:      "btn btn-default",
						type:           "button",
						"data-dismiss": "modal",
					}, "閉じる")
				)
			);
			break;

		case "success":
			content = m("div", {className: "modal-content"},
				m("div", {className: "modal-body"},
					primitive.successView),
				m("div", {className: "modal-footer"},
					m("button", {
						className:      "btn btn-default",
						type:           "button",
						"data-dismiss": "modal",
					}, "閉じる")
				)
			);
			break;

		default:
			throw new Error("unknown type: "+control.slides[0].type);
		}
	}

	return m("div",
		m("div", {
			ariaHidden:     control.loading || control.slides.length ? "true" : "false",
			ariaLabelledby: "component-password-modal-title",
			className:      "modal fade",
			config:         (function(element, initialized, context) {
				const jquery = $(element);

				if (!initialized) {
					context.onunload = jquery.modal.bind(jquery, "hide");

					jquery.on("hidden.bs.modal",
						() => !control.loading && !this.slides.length &&
							ensureRedraw(this.finish));
				}

				jquery.modal(control.loading || this.slides.length ? "hide" : "show");
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
						href:      "#!password",
						id:        "component-password-modal-title",
					}, primitive.titleView)
				), m("div", {className: "modal-body"},
					primitive.bodyView(control.primitive)
				), m("div", {className: "modal-footer"},
					primitive.buttonView(control.primitive)
				)
			)
		)),
		m("div", {
			ariaHidden: !control.loading || content ? "true" : "false",
			className:  "modal fade",
			config:     (function(element, initialized, context) {
				const jquery = $(element);

				if (!initialized) {
					context.onunload = jquery.modal.bind(jquery, "hide");
				}

				jquery.modal(!this.loading || this.slides.length ? "hide" : "show");
			}).bind(control),
			"data-backdrop": "static",
		}, m("div", {className: "modal-dialog", role: "document"},
			m("div", {className: "modal-body"},
				primitive.inprogressView
			)
		)),
		content && m("div", {
			className: "modal fade",
			config:    (function(element, initialized, context) {
				const jquery = $(element);

				if (!initialized) {
					context.onunload = jquery.modal.bind(jquery, "hide");

					jquery.on("hidden.bs.modal",
						ensureRedraw.bind(undefined,
							this.finishSlide.bind(this)));
				}

				jquery.modal("show");
			}).bind(control),
		}, m("div", {className: "modal-dialog", role: "document"}, content))
	);
}

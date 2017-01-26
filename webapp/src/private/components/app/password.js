/**
	@file password.js implements the component of the password updating
	page.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/password */

/**
	module:private/components/app/password is a component to provide the
	feature to update the password of the user.
	@name module:private/components/app/password
	@type external:Mithril~Component
*/

import * as container from "../container";
import * as primitive from "../password/primitive";
import * as progress from "../progress";

export function controller() {
	const control = {state: {}};

	control.primitive = new primitive.controller((function(promise) {
		this.state = {
			attributes: {value: 0},
			type:       "inprogress",
		};

		promise.then(() => {
			this.state = {type: "success"};
		}, message => {
			this.state = {type: "error", message};
		}, event => {
			this.state.attributes = {
				max:   event.total,
				value: event.loaded,
			};
		});
	}).bind(control));

	return control;
}

export function view(control) {
	let alertView;

	switch (control.state.type) {
	case "error":
		alertView = m("div", {
			className: "alert alert-danger",
			role:      "alert",
			style:     {
				display:   "inline-block",
				textAlign: "left",
			},
		}, primitive.errorView(control.state.message));
		break;

	case "inprogress":
		alertView = m("div", {
			className: "alert alert-info",
			role:      "alert",
			style:     {
				display:   "inline-block",
				textAlign: "left",
			},
		}, primitive.inprogressView);
		break;

	default:
		throw new Error("unknown type: "+control.state.type);
	}

	return [
		control.state.type == "inprogress" ? m(progress, control.state.attributes) : null,
		m(container,
			m("div", {className: "container", style: {textAlign: "center"}},
				m("form", {
					ariaLabelledby: "component-app-password-title",
					style:          {
						display:   "inline-block",
						textAlign: "left",
					},
				},
					m("div", {
						ariaHidden: (!alertView).toString(),
						style:      {
							minHeight: "8rem",
							textAlign: "center",
						},
					}, alertView),
					m("h1", {id: "component-app-password-title"},
						primitive.titleView),
					primitive.bodyView(control.primitive),
					m("div", {style: {textAlign: "center"}},
						primitive.buttonView(control.primitive))
				)
			),
			control.state.type == "success" ? m("div", {
				className: "modal fade",

				config(element, initialized, context) {
					const jquery = $(element);

					if (!initialized) {
						context.onunload = jquery.modal.bind(jquery, "hide");

						jquery.on("hidden.bs.modal",
							history.length ?
								history.back.bind(history) :
								() => m.route(""));
					}

					jquery.modal("show");
				},

				role:     "dialog",
				tabindex: "-1",
			}, m("div", {className: "modal-dialog", role: "document"},
				m("div", {className: "modal-content"},
					m("div", {className: "modal-body"},
						primitive.successView),
					m("div", {className: "modal-footer"},
						history.length ? m("button", {
							className:      "btn btn-default",
							type:           "button",
							"data-dismiss": "modal",
						}, "戻る") : m("a", {
							className: "btn btn-default",
							href:      "#",
						}, "トップページへ")
					)
				)
			)) : null
		),
	];
}

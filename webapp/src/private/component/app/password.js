/**
	@file password.js implements password component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/password */

/**
	module:private/component/app/password is a component to update the
	password of the user.
	@name module:private/component/app/password
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as primitive from "../password/primitive";

export function oninit() {
	this.primitive = primitive.newState();
}

export function view() {
	return m(container,
		m("div", {className: "container", style: {textAlign: "center"}},
			m("form", {
				"aria-labelledby": "component-app-password-title",
				style:             {
					display:   "inline-block",
					textAlign: "left",
				},
			},
				m("h1", {id: "component-app-password-title"},
					primitive.title),
				m(this.primitive.body, {
					autofocus: true,
					oncreate:  () => setTimeout(
						this.primitive.focus.bind(this.primitive)),
				}),
				m("div", {style: {textAlign: "center"}},
					m(this.primitive.button)))));
}

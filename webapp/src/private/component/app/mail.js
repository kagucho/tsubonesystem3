/**
	@file mail.js implements mail component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/mail */

/**
	module:private/component/app/mail is a component to mail.
	@name module:private/component/app/mail
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as primitive from "../mail/primitive";

export function oninit() {
	this.primitive = new primitive.default(m.route.param("subject"));
}

export function oncreate() {
	this.primitive.start();
}

export function onbeforeremove() {
	this.primitive.end();
}

export function view() {
	return m(container, {style: {height: "100%"}},
		m("form", {
			className: "container",
			style:     {
				display:       "flex",
				flexDirection: "column",
				height:        "100%",
			},
		},
			m("h1", primitive.title),
			m("div", {style: {flex: "1"}}, m(this.primitive.body)),
			m(this.primitive.button)));
}

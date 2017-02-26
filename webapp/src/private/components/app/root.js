/**
	@file index.js implements root component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/root */

/**
	module:private/components/app/root is a component to show the top page.
	@name module:private/components/app/root
	@type !external:Mithril~Component
*/

import * as container from "../container";
import Crossfade from "../../../crossfade";
import slides from "../../../slides";

export const view = () => m(container, m("div", {
	oncreate(node) {
		this.crossfade = new Crossfade(node.dom);
		this.crossfade.start();
	},

	onbreforeremove() {
		this.crossfade.stop();
	},

	style: {height: "100%", width: "100%"},
}, slides.map(image => m("div", {
	className: "crossfade-element",
	style:     {
		backgroundImage: "url(" + image + ")",
		backgroundSize:  "cover",
		position:        "fixed",
		height:          "100%",
		width:           "100%",
	},
}))));

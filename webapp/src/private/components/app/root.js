/**
	@file index.js implements the top page in the private session.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module private/components/app/root */

/**
	module:private/components/app/root is a component to provide the feature
	to show the top page.
	@name module:private/components/app/root
	@type external:Mithril~Component
*/

import * as container from "../container";
import Crossfade from "../../../crossfade";
import slides from "../../../slides";

export function view() {
	return m(container, m("div", {
		config(element, initialized, context) {
			if (!initialized) {
				const crossfade = new Crossfade(element);
				crossfade.start();
				context.onunload = crossfade.stop.bind(crossfade);
			}
		}, style: {height: "100%", width: "100%"},
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
}

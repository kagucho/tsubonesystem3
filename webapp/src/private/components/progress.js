/**
	@file progress.js implements the display of progress.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module private/components/progress */

/**
	module:private/components/progress is a component to show the progress.
	@name module:private/components/progress
	@type external:Mithril~Component
*/

export function view(control, attributes) {
	const max = attributes.max == null ? 1 : attributes.max;
	const value = attributes.value;

	const style = {
		backgroundColor: "lightgray",
		position:        "fixed",
		top:             "0px",
		height:          ".2rem",
		width:           "100%",
		zIndex:          "1031",
	};

	if (max == value) {
		style.animation = ".1s ease 1s 1 normal forwards running progress";
	}

	return m("div", {ariaHidden: "true", style}, m("div", {
		style: {
			backgroundColor: "dodgerblue", height:          "100%",
			width:           value / max * 100 + "%",
		},
	}));
}

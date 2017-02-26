/**
	@file progress.js implements the display of progress.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/progress */

/**
	module:private/components/progress is a component to show the progress.
	@name module:private/components/progress
	@type external:Mithril~Component
*/

export function view(node) {
	const max = node.attrs.max == null ? 1 : node.attrs.max;
	const {value} = node.attrs;

	return m("div", {
		"aria-hidden": "true",
		style:         {
			animation: max == value ?
				".1s ease 1s 1 normal forwards running component-progress-hiding" :
				"",
			backgroundColor: "lightgray",
			position:        "fixed",
			top:             "0",
			height:          "0.5rem",
			width:           "100%",

			/*
				its z-index should be larger than ones of bootstrap
				components.
				https://github.com/twbs/bootstrap/blob/v3.3.7/less/variables.less#L265
			*/
			zIndex: "1051",
		},
	}, m("div", {
		style: {
			backgroundColor: "dodgerblue",
			height:          "100%",
			transformOrigin: "0",
			transform:       `scaleX(${value / max})`,
		},
	}));
}

/**
	Attrs is a set of attributes for module:private/components/progress.
	@typedef module:private/components/progress~Attrs
	@property {?Number} max - The maximum value.
	@property {!Number} value - The current value.
	@see {@link https://html.spec.whatwg.org/multipage/forms.html#the-progress-element|
		HTML Standard 4.10.13 The progress element}
*/

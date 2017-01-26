/**
	@file container.js implements the container of the application.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/container */

/**
	module:private/components/container is a component to create a common
	container.
	@name module:private/components/container
	@type external:Mithril~Component
*/

export function view(control, ...children) {
	return m("div", [
		m("header", {
			className: "navbar navbar-default navbar-fixed-top",
			style:     {backgroundColor: "white", borderStyle: "hidden"},
		}, m("div", {className: "container"},
			m("div", {className: "navbar-header"},
				m("button", {
					ariaHidden:    "true",
					className:     "navbar-toggle",
					"data-toggle": "collapse",
					"data-target": ".navbar-top",
				}, m("span", {
					className: "glyphicon glyphicon-menu-hamburger",
				})),
				m("a", {
					className: "navbar-brand",
					href:      "#",
					style:     {color: "#2ca9e1"},
				}, "TsuboneSystem")
			),
			m("div", {
				className: "collapse navbar-collapse navbar-top",
				role:      "menu",
			},
				m("ul", {className: "nav navbar-nav"},
					m("li", m("a", {href: "#!mail"}, "Mail")),
					m("li", m("a", {href: "#!party"}, "Party")),
					m("li", m("a", {href: "#!members"}, "Members")),
					m("li", m("a", {href: "#!clubs"}, "Clubs")),
					m("li", m("a", {href: "#!officers"}, "Officers"))
				), m("ul", {className: "nav navbar-nav navbar-right"},
					m("li", m("a", {href: "#!signout"}, "Sign out")),
					m("li", m("a", {href: "#!settings"}, "Settings"))
				)
			)
		)), m("div", {style: {paddingTop: "50px"}}, ...children),
	]);
}

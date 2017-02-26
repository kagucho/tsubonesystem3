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

import client from "../client";
import {top} from "../navigator";

export function oninit() {
	client.userDetail().then(member => {
		this.nickname = member.nickname;
	});
}

export function view(node) {
	return m("div", {style: {height: "100%", position: "relative"}},
		m("header", {
			className: "navbar navbar-default navbar-fixed-top",
			style:     {backgroundColor: "white", borderStyle: "hidden"},
		}, m("nav", {className: "container"},
			m("div", {className: "navbar-header"},
				m("button", {
					"aria-hidden": "true",
					"data-toggle": "collapse",
					"data-target": ".navbar-top",
					className:     "navbar-toggle",
				}, m("span", {
					className: "glyphicon glyphicon-menu-hamburger",
				})),
				m("a", {
					className: "navbar-brand",
					href:      "",
					onclick:   top,
					style:     {color: "#2ca9e1"},
				}, "TsuboneSystem")
			),
			m("div", {
				className: "collapse navbar-collapse navbar-top",
				role:      "menu",
			},
				m("ul", {className: "nav navbar-nav"},
					m("li", m("a", {href: "#!mail"}, "Mail")),
					m("li", m("a", {href: "#!parties"}, "Parties")),
					m("li", m("a", {href: "#!members"}, "Members")),
					m("li", m("a", {href: "#!clubs"}, "Clubs")),
					m("li", m("a", {href: "#!officers"}, "Officers"))
				), m("ul", {className: "nav navbar-nav navbar-right"},
					m("li", m("a", {href: "#!signout"}, "Sign out")),
					this.nickname && m("li", m("a", {href: "#!member?id=" + this.id}, this.nickname))
				)
			)
		)), m("div", {
			style: {
				display:       "flex",
				flexDirection: "column",
				minHeight:     "100%",
			},
		},
			m("main", {style: {flex: "1", paddingTop: "50px"}},
				node.children
			), m("footer", {className: "clearfix"},
				m("div", {className: "container"},
					m("small",
						"Copyright © 2017 神楽坂一丁目通信局. 詳細は",
						m("a", {href: "/license"}, "こちら"),
						"をご覧ください。"
					)
				)
			)
		)
	);
}

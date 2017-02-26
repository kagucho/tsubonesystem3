/**
	@file mithril.js includes documentations for Mithril.
	This includes citations and you may refer to links for the source.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/**
	Mithril is a modern client-side Javascript framework for building Single
	Page Applications.
	@external Mithril
	@see {@link http://mithril.js.org/|Introduction - Mithril.js}
*/

/**
	Node represents an HTML element in a Mithril view.
	@interface external:Mithril~Node
	@extends Object
	@see {@link http://mithril.js.org/vnodes.html#structure|
		Virtual DOM nodes - Mithril.js Structure}
*/

/**
	external:Mithril~Children is a child vnode.
	@typedef {(external:Mithril~Node|String|Number|Boolean|external:Mithril~Children[])}
		external:Mithril~Children
*/

/**
	Component is an encapsulated part of a view to make code easier to
	organize and/or reuse.
	@typedef external:Mithril~Component
	@see {@link http://mithril.js.org/components.html|
		Components - Mithril.js}
*/

/**
	Routes is an object whose keys are route strings and values are either
	components or a RouteResolver.
	@typedef external:Mithril~Routes
	@see {@link http://mithril.js.org/route.html#signature|
		route(root, defaultRoute, routes) - Mithril.js}
*/

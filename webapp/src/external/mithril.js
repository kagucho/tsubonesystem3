/**
	@file mithril.js includes documentations for Mithril.
	This includes citations and you may refer to links for the source.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/**
	Mithril is a client-side MVC framework
	- a tool to organize code in a way that is easy to think about and to
	maintain.
	@external Mithril
	@see {@link http://mithril.js.org/mithril.html|Mithril API}
*/

/**
	m is a convenience method to compose virtual elements that can be
	rendered via m.render().
	@class external:Mithril.m
	@extends external:ES.Object
	@see {@link http://mithril.js.org/mithril.html|Mithril API}
*/

/**
	Children can be rendered via m.render().
	@typedef {(external:Mithril.m|external:ES.String|external:Mithril~Children[])}
		external:Mithril~Children
*/

/**
	Component is a building block for Mithril applications.
	They allow developers to encapsulate functionality into reusable units.
	@typedef external:Mithril~Component
	@type {external:ES.Object}
	@property {?external:Mithril~Controller} controller - The
	function which returns the control.
	@property {!external:Mithril~View} view - The function
	which returns the children.
	@see {@link http://mithril.js.org/mithril.component.html|
		m.component - Mithril}
*/

/**
	Controller returns the control.
	@callback external:Mithril~Controller
	@param {...*} - The arguments passed to m.component except the component
	itself.
	@returns {?external:ES.Object} The control.
*/

/**
	View returns the children.
	@callback external:Mithril~View
	@param {?external:ES.Object} - The control returned by
	external:Mithril~Controller.
	@param {...*} - The arguments passed to m.component except the component
	itself.
	@returns {?external:Mithril~Children} The children.
*/

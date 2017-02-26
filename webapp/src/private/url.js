/**
	@file url.js implements URL encoders.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
	@see {@link https://tools.ietf.org/html/rfc7595|
		RFC 7595 - Guidelines and Registration Procedures for URI Schemes}
*/

/** @module private/url */

/**
	mailto returns the URL of the given mail address.
	@param {!String} address - The mail address.
	@returns {!String} The corresponding URL.
	@see {@link https://tools.ietf.org/html/rfc6068|
		RFC 6068 - The 'mailto' URI Scheme}
*/
export const mailto = address => "mailto:" + punycode.toASCII(encodeURIComponent(address));

/**
	tel returns the URL of the given telphone number.
	@param {!String} number - The telphone number.
	@returns {!String} The corresponding URL.
	@see {@link https://tools.ietf.org/html/rfc3966|
		RFC 3966 - The tel URI for Telephone Numbers}
*/
export const tel = number => "tel:" + number.replace(/^0(?!-)/, "+81-");

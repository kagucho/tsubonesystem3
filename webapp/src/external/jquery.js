/**
	@file jquery.js includes documentations for jQuery.
	This includes citations and you may refer to links for the source.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/**
	jQuery is a fast, small, and feature-rich JavaScript library.
	@external jQuery
	@see {@link http://api.jquery.com/|jQuery API Documentation}
*/

/**
	jQuery object contains a collection of Document Object Model (DOM)
	elements that have been created from an HTML string or selected from a
	document.
	@interface external:jQuery~jQuery
	@extends Object
	@see {@link http://api.jquery.com/Types/#jQuery|
		Types | jQuery API Documentation jQuery}
*/

/**
	This object provides a subset of the methods of the Deferred object.
	It prevents users from changing the state of the Deferred.
	@interface external:jQuery~Promise
	@extends Object
	@see {@link http://api.jquery.com/Types/#Promise|
		Types | jQuery API Documentation Promise Object}
*/

/**
	always adds handlers to be called when the external:jQuery~Promise is
	either resolved or rejected.
	@function external:jQuery~Promise#always
	@param {!Function|Function[]} callback - A
	function, or array of functions, that is called when the Deferred is
	resolved or rejected.
	@param {...!(Function|Function[])}
	[additionalCallbacks] - Optional additional functions, or arrays of
	functions, that are called when the Deferred is resolved or rejected.
	@returns {!external:jQuery~Promise} The
	external:jQuery~Promise.
	@see {@link http://api.jquery.com/deferred.always/|
		deferred.always() | jQuery API Documentation}
*/

/**
	catch adds handlers to be called when the Deferred object is rejected.
	@function external:jQuery~Promise#catch
	@param {!Function} filter - A function that is called when
	the Deferred is rejected.
	@returns {!external:jQuery~Promise} A filtered promise.
	@see {@link http://api.jquery.com/deferred.catch/|
		deferred.catch() | jQuery API Documentation}
*/

/**
	done adds handlers to be called when the Deferred object is resolved.
	@function external:jQuery~Promise#done
	@param {!Function|Function[]} callback - A
	function, or array of functions, that are called when the Deferred is
	resolved.
	@param {...!(Function|Function[])}
	[additionalCallbacks] - Optional additional functions, or arrays of
	functions, that are called when the Deferred is resolved.
	@returns {!external:jQuery~Promise} The
	external:jQuery~Promise.
	@see {@link http://api.jquery.com/deferred.done/|
		deferred.done() | jQuery API Documentation}
*/

/**
	progress adds handlers to be called when the Deferred object generates
	progress notifications.
	@function external:jQuery~Promise#progress
	@param {!Function|Function[]} callback - A
	function, or array of functions, to be called when the Deferred
	generates progress notifications.
	@param {...!(Function|Function[])}
	[additionalCallbacks] - Optional additional functions, or arrays of
	functions, to be called when the Deferred generates progress
	notifications.
	@returns {!external:jQuery~Promise} The
	external:jQuery~Promise.
	@see {@link http://api.jquery.com/deferred.progress/|
		deferred.progress() | jQuery API Documentation}
*/

/**
	promise returns external:jQuery~Promise itself.
	@function external:jQuery~Promise#promise
	@param {?*} target - Object onto which the promise methods have to be
	attached.
	@returns {!external:jQuery~Promise} itself, or target
	if provided.
	@see {@link http://api.jquery.com/deferred.promise/|
		deferred.promise() | jQuery API Documentation}
*/

/**
	state determines the current state of a Deferred object.
	@function external:jQuery~Promise#state
	@returns {!String} A string representing the current state
	of the Deferred object.
	@see {@link http://api.jquery.com/deferred.state/|
		deferred.state() | jQuery API Documentation}
*/

/**
	then adds handlers to be called when the Deferred object is resolved,
	rejected, or still in progress.
	@function external:jQuery~Promise#then
	@param {?*} onresolved -  A function, or array of functions, called
	when the Deferred is resolved.
	@param {?*} [onrejected] - A function, or array of functions, called
	when the Deferred is rejected.
	@param {?*} [onprogress] - A function, or array of functions, called
	when the Deferred notifies progress.
	@returns {!external:jQuery~Promise} A filtered promise.
	@see {@link http://api.jquery.com/deferred.then/|
		deferred.then() | jQuery API Documentation}
*/

/**
	The jqXHR object.
	@interface external:jQuery~jqXHR
	@extends external:jQuery~Promise
	@see {@link http://api.jquery.com/jQuery.ajax/#jqXHR|
		jQuery.ajax() | jQuery API Documentation}
*/

/**
	Event is an object guaranteed to be passed to the event handler. Most
	properties from the original event are copied over and normalized to the
	new event object.
	@typedef external:jQuery~Event
	@see {@link http://api.jquery.com/category/events/event-object/|
		Event Object | jQuery API Documentation}
*/

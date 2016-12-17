/**
 * @file jquery.js includes documentations for jQuery.
 * This includes citations and you may refer to links for the source.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/**
 * jQuery is a fast, small, and feature-rich JavaScript library.
 * @external jQuery
 * @see {@link http://api.jquery.com/|jQuery API Documentation}
 */

/**
 * jQuery API.
 * @class external:jQuery.$
 * @extends external:ES.Object
 * @see {@link http://api.jquery.com/|jQuery API Documentation}
 */

/**
 * A chainable utility object created by calling the external:jQuery#$.Deferred
 * method.
 * @class external:jQuery.$.Deferred
 * @extends external:jQuery.$.Deferred#promise
 * @param {?external:jQuery.$.Deferred~beforeStart} beforeStart - A function
 * that is called just before the constructor returns.
 * @returns {!external:jQuery.$.Deferred} - A new instance.
 * @see {@link http://api.jquery.com/jQuery.Deferred/|
               jQuery.Deferred() | jQuery API Documentation}
 */

/**
 * This object provides a subset of the methods of the Deferred object.
 * It prevents users from changing the state of the Deferred.
 * @class external:jQuery.$.Deferred#promise
 * @extends external:ES.Object
 * @param {?external:ES.Object} target - Object onto which the promise methods
 * have to be attached.
 * @returns {!external:jQuery.$.Deferred#promise} - A new instance.
 * @see {@link http://api.jquery.com/deferred.promise/|
               deferred.promise() | jQuery API Documentation}
 */

/**
 * A function that is called just before external:jQuery#$.Deferred returns.
 * @param {!external:jQuery.$.Deferred} deferred - A new instance.
 * @callback external:jQuery.$.Deferred~beforeStart
 */

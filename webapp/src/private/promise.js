/**
	@file promise.js implements a wrapper for jQuery Promise.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/promise */

/**
	stateMap is a WeakMap which holds states.
	@private
	@type !WeakMap<!module:private/promise~State>
*/
const stateMap = new WeakMap;

/**
	redrawOpportunity gives Mithril an opportunity to redraw.
	@private
	@this module:private/promise~State
	@param {!String} key - The key to identify the pseudo event.
	@returns {Undefined}
*/
function redrawOpportunity(key) {
	try {
		if (this[key].redraw !== false) {
			m.redraw();
		}
	} catch (exception) {
		debugger;
	}
}

/**
	onalwaysWrapper is a function to wrap a callback to be called when
	a external:jQuery~Promise is resolved or rejected.
	@private
	@this module:private/promise~State
	@param {!function} callback - The callback to be wrapped.
	@param {...*} wrapped - The arguments to be passed to the callback.
	@returns {*} The value returned by the callback.
*/
function onalwaysWrapper(callback, ...wrapped) {
	try {
		const result = callback(...wrapped);

		if (result && $.isFunction(result.state)) {
			const state = stateMap.get(result);
			if (state) {
				delete this.always;
				delete this.progress;
				Object.setPrototypeOf(this, state);

				return state.promise;
			}

			if (result.state() == "pending") {
				const event = {redraw: false};

				this.always = event;
				this.progress = event;

				return result;
			}
		}

		return result;
	} finally {
		redrawOpportunity.call(Object.getPrototypeOf(this), "always");
	}
}

/**
	wrapChainer returns a wrapped chaining function of
	external:jQuery~Promise.
	@private
	@param {!module:private/promise} promise - The wrapping promise.
	@param {!module:private/promise~Chainer} chainer - The chaining function
	to be wrapped.
	@param {!String} event - The pseudo event to be associated.
	@param {!Iterable<!(function|function[])>}
	callbacks - The callbacks to be associated.
	@returns {!module:private/promise} The wrapping promise.
*/
function wrapChainer(promise, chainer, event, callbacks) {
	const state = stateMap.get(promise);
	state.promise[chainer](...callbacks, redrawOpportunity.bind(state, event));

	return promise;
}

/**
	eventProxy returns a new proxy to capture redraw property of event.
	@private
	@param {*} event - The target.
	@param {!module:private/promise~State} state - The state to contain a
	pseudo event to capture redraw property.
	@param {!String} key - The key to identify the pseudo event
	to capture redraw property.
	@returns {*} The proxy.
*/
function eventProxy(event, state, key) {
	state[key] = {};

	return new Proxy(event, {
		get: function(target, targetKey) {
			return (targetKey == "redraw" ? this : target)[targetKey];
		}.bind(state[key]),

		set: function(target, targetKey, value) {
			return (targetKey == "redraw" ? this : target)[targetKey] = value;
		}.bind(state[key]),
	});
}

/**
	module:private/promise~Promise is a class wrapping
	external:jQuery~Promise.
	@implements {external:jQuery~Promise}
*/
class Promise {
	/**
		constructor constructs a new module:private/promise.
		@param {!module:private/promise~State} state - The state.
		@returns {Undefined}
	*/
	constructor(state) {
		stateMap.set(this, state);
		Object.freeze(this);
	}

	always() {
		return wrapChainer(this, "always", "always", arguments);
	}

	catch(onerror) {
		const state = stateMap.get(this);
		const newState = Object.create(state);

		newState.promise = state.promise.catch(onalwaysWrapper.bind(newState, onerror));

		return new this.constructor(newState);
	}

	done() {
		return wrapChainer(this, "done", "always", arguments);
	}

	progress() {
		return wrapChainer(this, "progress", "progress", arguments);
	}

	promise(target) {
		if (target == null) {
			return this;
		}

		stateMap.set($.extend(target, this), stateMap.get(this));

		return target;
	}

	state() {
		return stateMap.get(this).promise.state();
	}

	then(onsuccess, onerror, onprogress) {
		const state = stateMap.get(this);

		if ($.isArray(onsuccess) || $.isArray(onerror) || $.isArray(onprogress)) {
			const alwaysRedrawOpportunity = redrawOpportunity.bind(state, "always");
			const progressRedrawOpportunity = redrawOpportunity.bind(state, "progress");

			return state.promise.then(
				$.isArray(onsuccess) ?
					onsuccess.concat(alwaysRedrawOpportunity) :
					($.isFunction(onsuccess) ?
						[onsuccess, alwaysRedrawOpportunity] :
						onsuccess, onerror, onprogress),
				$.isArray(onerror) ?
					onerror.concat(alwaysRedrawOpportunity) :
					($.isFunction(onerror) ?
						[onerror, alwaysRedrawOpportunity] :
						onerror),
				$.isArray(onprogress) ?
					onprogress.concat(progressRedrawOpportunity) :
					($.isFunction(onprogress) ?
						[onprogress, progressRedrawOpportunity] :
						onprogress));
		}

		const newState = Object.create(state);

		newState.promise = state.promise.then(
			$.isFunction(onsuccess) && onalwaysWrapper.bind(newState, onsuccess),
			$.isFunction(onerror) && onalwaysWrapper.bind(newState, onerror),
			$.isFunction(onprogress) && function() {
				const prototype = Object.getPrototypeOf(newState);

				try {
					return onprogress(...arguments);
				} catch (exception) {
					newState.always = prototype.progress;
					throw exception;
				} finally {
					redrawOpportunity.call(prototype, "progress");
				}
			}
		);

		return new this.constructor(newState);
	}
}

/**
	module:private/promise.when wraps jQuery.when().
	@param {...*} [values] - The values to be reflected to the resolved
	value.
	@returns {!module:private/promise} - A promise.
	@see {@link http://api.jquery.com/jQuery.when/|jQuery.when() | jQuery API Documentation}
*/
export function when() {
	const result = Array(arguments.length);
	const newState = {promise: $.Deferred()};
	let promises = arguments.length;

	Array.prototype.forEach.call(arguments, (value, index) => {
		if (value && $.isFunction(value.then)) {
			const state = stateMap.get(value);

			if (state) {
				state.promise.then(function() {
					result[index] = arguments.length > 1 ? [...arguments] : arguments[0];

					promises--;
					if (promises) {
						return;
					}

					newState.always = state.always;
					newState.promise.resolve(...result);
				}, function() {
					Object.setPrototypeOf(newState, state);
					newState.promise.reject(...arguments);
				});
			} else {
				value.then(function() {
					result[index] = arguments.length > 1 ? [...arguments] : arguments[0];

					promises--;
					if (!promises) {
						newState.promise.resolve(...result);
					}
				}, function() {
					newState.prmise.reject(...arguments);
				});
			}
		} else {
			result[index] = value;
			promises--;
		}
	});

	return promises ?
		new Promise(newState) : newState.promise.resolve(result);
}

/**
	wrap returns a wrapped promise which simply redraws after any calls of
	callbacks.
	@param {!external:jQuery~Promise} promise - A promise to be wrapped.
	@returns {!module:private/promise~Promise} A wrapped promise.
*/
export const wrap =
	promise => new Promise({always: {}, progress: {}, promise});

/**
	module:private/promise is a class to wrap external:jQuery~Promise.
	@extends Object
*/
export class Wrapper {
	/**
		constructor constructs a module:private/promise.
		@returns {Undefined}
	*/
	constructor() {
		stateMap.set(this, {});
		Object.freeze(this);
	}

	/**
		alwaysProxy returns a proxy for events of the resolution
		or the rejection.
		@param {*} event - The event.
		@returns {*} The proxy event.
	*/
	alwaysProxy(event) {
		return eventProxy(event, stateMap.get(this), "always");
	}

	/**
		progressProxy returns a proxy for events of the notification.
		@param {*} event - The event.
		@returns {*} The proxy event.
	*/
	progressProxy(event) {
		return eventProxy(event, stateMap.get(this), "progress");
	}

	/**
		wrap Wraps the given external:jQuery~Promise.
		It can be called only once.
		@param {!external:jQuery~Promise} promise - The promise to be
		wrapped.
		@returns {!module:private/promise~Promise} - The wrapped
		promise.
	*/
	wrap(promise) {
		const state = stateMap.get(this);

		if (state.promise) {
			throw new Error("promise is already set");
		}

		state.promise = promise;

		return new Promise(state);
	}
}

/**
	A function of a promise which can be chained.
	@private
	@callback module:private/promise~Chainer
	@params {...function} callbacks - The callbacks.
	@returns {!external:jQuery~Promise} The promise.
*/

/**
	A pseudo event of module:private/promise to capture redraw property/
	@private
	@typedef module:private/promise~Event
	@extends {Object}
	@property {?Boolean} redraw - Mithril will redraw if redraw
	property is NOT STRICTLY false.
*/

/**
	A state of module:private/promise.
	@private
	@typedef module:private/promise~State
	@extends {Object}
	@property {!external:jQuery~Promise} promise - The wrapped promise.
	@property {?module:private/promise~Event} always - The pseudo event
	which will get fired when the promise gets resolved or rejected.
	@property {?module:private/promise~Event} progress - The pseudo event
	which will get fired when the promise gets notified.
*/

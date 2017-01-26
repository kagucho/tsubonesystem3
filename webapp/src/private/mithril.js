/**
	@file mithril.js implements utilities for Mithril.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/mithril */

/**
	ensureRedraw ensures Mithril redraws.
	@param {!external:ES.Function|external:jQuery.$.Deferred#promise}
	object - a callback which may or may not redraw, or a promise which
	may not ensure to redraw after resolved, rejected, or notified.
	@returns {external:ES~Undefined|external:jQuery.$.Deferred#promise}
	If the object is a promise, it returns a promise which resolves,
	rejects, or notifies as the given promise does so, and ensures the
	callbacks redraws.
*/
export function ensureRedraw(object) {
	if ($.isFunction(object)) {
		const {redraw} = m;

		m.redraw = force => {
			if (force) {
				if (redraw.strategy() != "none") {
					m.redraw = redraw;
					m.redraw(force);
				}
			}
		};

		m.redraw.strategy = redraw.strategy;

		object();

		if (m.redraw != redraw) {
			m.redraw = redraw;
			m.redraw();
		}
	} else {
		const wrapped = Object.create(object);

		for (const key of ["always", "catch", "done", "progress", "then"]) {
			wrapped[key] = (function(...callbacks) {
				return ensureRedraw(this(...(function *() {
					for (const callback of callbacks) {
						yield (function() {
							ensureRedraw(() => this(...arguments));
						}).bind(callback);
					}
				}())));
			}).bind(object[key].bind(object));
		}

		return wrapped;
	}
}

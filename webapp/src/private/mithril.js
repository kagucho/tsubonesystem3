/**
	@file mithril.js implements utilities for Mithril.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/mithril */

let ensuring;

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
		if (ensuring) {
			object();

			return;
		}

		const {redraw} = m;

		try {
			ensuring = true;

			m.redraw = force => {
				if (redraw.strategy() == "none") {
					redraw.strategy("diff");
				} else if (force) {
					m.redraw = redraw;
					m.redraw(force);
				}
			};

			m.redraw.strategy = redraw.strategy;

			object();
		} finally {
			ensuring = false;

			if (m.redraw != redraw) {
				m.redraw = redraw;
				m.redraw();
			}
		}
	} else {
		const wrapped = Object.create(object);

		for (const key of ["always", "catch", "done", "progress", "then"]) {
			wrapped[key] = (function(...callbacks) {
				return ensureRedraw(this(...(function *(boundCallbacks) {
					for (const callback of boundCallbacks) {
						yield (function() {
							let result;

							ensureRedraw(() => {
								result = this(...arguments);
							});

							return result;
						}).bind(callback);
					}
				}(callbacks))));
			}).bind(object[key].bind(object));
		}

		return wrapped;
	}
}

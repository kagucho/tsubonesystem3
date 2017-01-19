/**
	@file mithril.js implements utilities for Mithril.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module private/mithril */

/**
	ensureRedraw ensures Mithril redraws.
	@param {!external:ES.Function} callback - The callback which may or may not redraw.
	@returns {external:ES~Undefined}
*/
export function ensureRedraw(callback) {
	const redraw = m.redraw;

	m.redraw = force => {
		if (force) {
			if (redraw.strategy() != "none") {
				m.redraw = redraw;
				m.redraw(force);
			}
		}
	};

	m.redraw.strategy = redraw.strategy;

	callback();

	if (m.redraw != redraw) {
		m.redraw = redraw;

		if (redraw.strategy() != "none") {
			m.redraw();
		}
	}
}

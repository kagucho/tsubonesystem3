/**
	@file navigator.js implements a feature to navigate.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/navigator */

/**
	initLength is the initial history.length.
	@private
	@type !Number
*/
const initLength = history.length;

/**
	back backs the history.
	@returns {Undefined}
*/
export function back() {
	history.back();
}

/**
	top navigates to the top page.
	@returns {Undefined}
*/
export function top() {
	history.pushState(null, "", location.pathname);
	dispatchEvent(new PopStateEvent("popstate", {bubbles: true}));
}

/**
	leaver returns an appropriate function to leave the current page. The
	returned value is either back or top.
	@function
	@returns {function} A function to leave.
*/
export const leaver = () => history.length > initLength ? back : top;

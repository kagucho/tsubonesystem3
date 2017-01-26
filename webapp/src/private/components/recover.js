/**
	@file recover.js implements the feature to recover the session.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/recover */

/**
	module:private/components/recover is a component to recover the session.
	@name module:private/components/recover
	@type external:Mithril~Component
*/

import * as progress from "./progress";
import {recoverSession} from "../client";

/**
	boundProgress holds the progress bound to the recovery.
	@type !external:ES.WeakMap<external:jQuery.$.Deferred#promise, external:ES.Number>
*/
const boundProgress = new WeakMap;

export function controller() {
	const promise = recoverSession();
	boundProgress.set(promise, {value: 0});

	promise.progress(function(event) {
		boundProgress.set(this, {max: event.total, value: event.loaded});
	}.bind(promise));

	return promise;
}

export function view(promise) {
	const current = boundProgress.get(promise);

	return current && m(progress, current);
}

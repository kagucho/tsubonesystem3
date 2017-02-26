/**
	@file index.js implements the entry point for the private page.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private */

import "./navigator";

/**
	container is the element which will contain the content.
	@private
	@type {!external:DOM~HTMLElement}
*/
const container = document.getElementById("container");

/**
	client is module:private/client
	@private
	@type ?module:private/client
*/
let client;

/**
	session is a promise which will get resolved or rejected when the
	recovery of the session was succeeded or failed.
	@private
	@type ?module:private/promise
*/
let session;

/**
	sessionProgress is the set of attributes for
	module:private/components/progress.
	@private
	@type !module:private/components/progress~Attrs
*/
let sessionProgress = {value: 0};

/**
	pendings are pending executions depending on asynchronously executing
	scripts.
	@private
	@type !module:private~Pending[]
*/
const pendings = [
	{
		dependencies: {deferreds: true, polyfill: true},
		callback() {
			client = require("./client").default;
			session = client.recoverSession().progress(event => {
				sessionProgress = {
					max:   event.total,
					value: event.loaded,
				};
			});

			sessionProgress = {value: 0};
		},
	}, {
		dependencies: {
			deferreds: true, mithril:   true,
			polyfill:  true, punycode:  true,
		},
		callback() {
			const app = require("./components/app");
			const query = m.parseQueryString(
				location.hash.slice(location.hash.indexOf("?")));

			if (query.fill) {
				client.setFillingToken(query.id, query.fill);
				m.route(container, "", app.fill);
			} else {
				m.mount(container, {
					view() {
						return m(require("./components/progress"),
							sessionProgress);
					},
				});

				session.catch(() => {
					const deferred = $.Deferred();

					m.mount(container, {
						view() {
							return m(require("./components/signin"), {
								onloadstart(promise) {
									promise.done(deferred.resolve.bind(deferred));
								},
							});
						},
					});

					return deferred;
				}).done(m.route.bind(m, container, "", app.default));
			}
		},
	},
];

/**
	execute executes pending executions according to the resolved
	dependency.
	@function execute
	@global
	@param {!String} dependency - The ID of the resolved dependency.
*/
window.execute = pendings.forEach.bind(pendings, function(pending, index) {
	delete pending.dependencies[this];

	for (const key in pending.dependencies) {
		return;
	}

	pending.callback();
	delete pendings[index];
});

onload = () => window.execute("deferreds");

/**
	@typedef module:private~Pending
	@private
	@property {Object.<String, *>} dependencies -
	An object which have enumerable properties whose keys are the names of
	the dependencies.
	@property {function} callback - The callback to be called
	when all of the dependencies are resolved.
*/

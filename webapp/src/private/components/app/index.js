/**
	@file index.js implements the route of the application.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app */

import * as club from "./club";
import * as clubs from "./clubs";
import * as mail from "./mail";
import * as member from "./member";
import * as members from "./members";
import * as notfound from "./notfound";
import * as notfoundFilling from "./notfound_filling";
import * as officer from "./officer";
import * as officers from "./officers";
import * as parties from "./parties";
import * as party from "./party";
import * as password from "./password";
import * as root from "./root";

/*
	FIX UPSTREAM: The fallback for "not found" pages are implemented with
	the order of properties. This should be fixed by altering Mithril API.
*/

/**
	module:private/components/app.normal is the internal representation of
	module:private/components/app
	@private
	@type !external:Mithril~Routes
*/
const normal = {
	"":          root,
	club, clubs, mail, member, members, officer, officers, party, parties,
	password,
	":route...": notfound,
};

/**
	module:private/components/app.fill is the routes in the filling mode,
	when the user agent is prompting the user to fill his information.
	@type !external:Mithril~Routes
*/
export const fill = {
	"":          root,
	member,
	":route...": notfoundFilling,
};

/**
	module:private/components/app is the routes in the normal mode.
	@type !external:Mithril~Routes
*/
export default normal;

for (const routes of [normal, fill]) {
	Object.freeze(routes);

	for (const route in routes) {
		Object.freeze(route);
	}
}

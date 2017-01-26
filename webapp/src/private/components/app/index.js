/**
	@file index.js implements the route of the application.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app */

import * as club from "./club";
import * as clubs from "./clubs";
import * as member from "./member";
import * as members from "./members";
import * as officer from "./officer";
import * as officers from "./officers";
import * as password from "./password";
import * as root from "./root";

/**
	module:private/components/app is the route of the application.
	@type external:ES.Object<external:ES.String, external:Mithril~Component>
*/
export default Object.freeze({
	"":          root,
	"!club":     club, "!clubs":    clubs,
	"!member":   member, "!members":  members,
	"!officer":  officer, "!officers": officers,
	"!password": password,
});

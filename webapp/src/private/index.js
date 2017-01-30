/**
	@file index.js implements the entry point for the private page.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private */

import * as recover from "./components/recover";
import * as signin from "./components/signin";
import app from "./components/app";

/**
	container is the element which will contain the content.
	@type {!external:DOM~HTMLElement}
*/
const container = document.getElementById("container");

m.mount(container, recover).catch(
	m.mount.bind(m, container, signin)).done(
		() => m.route(container, "", app));

m.route.mode = "hash";

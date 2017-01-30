/**
	@file index.js implements the entry point for the public page.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module public */

import Crossfade from "./crossfade";

document.getElementById("about-button").onclick =
	() => $(document.scrollingElement).animate(
		{scrollTop: document.getElementById("about").offsetTop});

(new Crossfade(document.getElementById("background"))).start();

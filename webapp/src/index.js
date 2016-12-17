/**
 * @file index.js implements the entry point for the public page.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/** @module public */

import Crossfade from "./crossfade";

/**
 * startCrossfade starts crossfade animation of background.
 * @returns {external:ES~Undefined}
 */
function startCrossfade() {
  const instance = new Crossfade(document.getElementById("background"));
  instance.start();
}

document.getElementById("about-button").onclick = () => {
  $(document.scrollingElement).animate(
    {scrollTop: document.getElementById("about").offsetTop});
};

startCrossfade();

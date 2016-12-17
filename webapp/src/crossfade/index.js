/**
 * @file index.js is the main file of crossfade module.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/** @module crossfade */

/*
import "./index.css";
*/

/** module:crossfade is a class which implements crossfade animation.
 * @extends external:ES.Object
 */
export default class {
  /**
   * constructor constructs a new instance.
   * @param {!external:DOM~Element} container - An element which contains
   * elements to crossfade.
   */
  constructor(container) {
    this.container = container;
  }

  /**
   * start starts crossfade animation.
   * @returns {external:ES~Undefined}
   */
  start() {
    let next = 0;

    const swap = () => {
      let previous = next - 2;
      if (previous < 0)
        previous += this.container.children.length;

      let current = next - 1;
      if (current < 0)
        current += this.container.children.length;

      this.container.children[previous].className = "crossfade-element";
      this.container.children[current].className = "crossfade-visible";
      this.container.children[next].className = "crossfade-fadingin";

      next++;
      if (next >= this.container.children.length)
        next = 0;
    };

    swap();
    this.interval = setInterval(swap, 8000);
  }

  /**
   * stop stops crossfade animation.
   * @returns {external:ES~Undefined}
   */
  stop() {
    clearInterval(this.interval);
  }
}

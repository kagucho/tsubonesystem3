/**
 * @file webpack.js includes documentations for Webpack.
 * This includes citations and you may refer to links for the source.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/**
 * Webpack is a module bundler.
 * @external Webpack
 */

/**
 * Loaders allow you to preprocess files as you require() or “load” them.
 * @interface external:Webpack~Loader
 * @see {@link https://webpack.github.io/docs/loaders.html|loaders}
 */

/**
 * @function external:Webpack~Loader.pitch
 * @param {!external:ES.String} remainingRequest - The remaining request.
 * @param {!external:ES.String} precedingRequest - The preceding request.
 * @param {!external:ES.Object} data - The presistent context.
 * @returns {?external:ES.String} The transformed content.
 * @see {@link https://webpack.github.io/docs/loaders.html#pitching-loader|
 *       loaders pitching loader}
 */

/**
 * @file es.js includes documentations for ECMAScript.
 * This includes citations and you may refer to links for the source.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/**
 * ECMAScript is the language embedded in web browsers.
 * @external ES
 * @see {@link
 *       http://www.ecma-international.org/ecma-262/6.0/#sec-object-objects|
 *       ECMAScript 2015 Language Specification – ECMA-262 6th Edition}
 */

/**
 * The Undefined type has exactly one value, called undefined.
 * @interface external:ES~Undefined
 * @extends external:ES.Object
 * @see {@link
 *       http://www.ecma-international.org/ecma-262/6.0/#sec-ecmascript-language-types-undefined-type|
 *       ECMAScript 2015 Language Specification – ECMA-262 6th Edition
 *       6.1.1 The Undefined Type}
 */

/**
 * Object object.
 * @class external:ES.Object
 * @param {?external:ES.Object} value - The value to be converted to a value of
 * type Object.
 * @returns {!external:ES.Object} A new object.
 * @see {@link
 *       http://www.ecma-international.org/ecma-262/6.0/#sec-object-objects|
 *       ECMAScript 2015 Language Specification – ECMA-262 6th Edition
 *       19.1 Object Objects}
 */

/**
 * String object.
 * @class external:ES.String
 * @extends external:ES.Object
 * @param {?external:ES.Object} value - The value to be converted to String.
 * @returns {!external:ES.String} A new string.
 * @see {@link
 *       http://www.ecma-international.org/ecma-262/6.0/#sec-string-objects|
 *       ECMAScript 2015 Language Specification – ECMA-262 6th Edition
 *       21.1 String Objects}
 */

/**
 * WeakMap objects are collections of key/value pairs where the keys are objects
 * and values may be arbitrary ECMAScript language values.
 * @class external:ES.WeakMap
 * @extends external:ES.Object
 * @param {?external:ES~Iterable} iterable - The value to be converted to
 * WeakMap.
 * @returns {!external:ES.WeakMap} A new WeakMap.
 * @see {@link
 *       http://www.ecma-international.org/ecma-262/6.0/#sec-weakmap-objects|
 *       ECMAScript 2015 Language Specification – ECMA-262 6th Edition
 *       23.3 WeakMap Objects}
 */

/**
 * The Iterable interface includes @@iterator property.
 * @interface external:ES~Iterable
 * @see {@link
 *       http://www.ecma-international.org/ecma-262/6.0/#sec-iterable-interface|
 *       ECMAScript 2015 Language Specification – ECMA-262 6th Edition
 *       25.1.1.1 The Iterable Interface}
 */

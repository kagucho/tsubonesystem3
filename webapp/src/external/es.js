/**
	@file es.js includes documentations for ECMAScript.
	This includes citations and you may refer to links for the source.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/**
	ECMAScript is the language embedded in web browsers.
	@external ES
	@see {@link
		http://www.ecma-international.org/ecma-262/7.0/index.html|
		ECMAScript® 2016 Language Specification}
*/

/**
	The type whose sole value is the undefined value.
	@interface external:ES~Undefined
	@extends external:ES.Object
	@see {@link
		http://www.ecma-international.org/ecma-262/7.0/index.html#sec-terms-and-definitions-undefined-type|
		ECMAScript® 2016 Language Specification
		4.3.11 Undefined type}
*/

/**
	The type consisting of the primitive values true and false.
	@class external:ES.Boolean
	@extends external:ES.Object
	@see {@link
		http://www.ecma-international.org/ecma-262/7.0/index.html#sec-terms-and-definitions-boolean-type|
		ECMAScript® 2016 Language Specification
		4.3.15 Boolean type}
*/

/**
	The set of all possible Number values including the special
	“Not-a-Number” (NaN) value, positive infinity, and negative infinity.
	@class external:ES.Number
	@see {@link
		http://www.ecma-international.org/ecma-262/7.0/index.html#sec-terms-and-definitions-number-type|
		ECMAScript® 2016 Language Specification
		4.3.21 Number type}
*/

/**
	The set of all possible String values.
	@class external:ES.String
	@extends external:ES.Object
	@see {@link
		http://www.ecma-international.org/ecma-262/7.0/index.html#sec-terms-and-definitions-string-type|
		ECMAScript® 2016 Language Specification
		4.3.18 String type}
*/

/**
	Object object.
	@class external:ES.Object
	@see {@link
		http://www.ecma-international.org/ecma-262/7.0/index.html#sec-object-objects|
		ECMAScript® 2016 Language Specification
		19.1 Object Objects}
*/

/**
	WeakMap objects are collections of key/value pairs where the keys are objects
	and values may be arbitrary ECMAScript language values.
	@class external:ES.WeakMap
	@extends external:ES.Object
	@see {@link
		http://www.ecma-international.org/ecma-262/7.0/index.html#sec-weakmap-objects|
		ECMAScript® 2016 Language Specification
		23.3 WeakMap Objects}
*/

/**
	The Iterable interface includes @@iterator property.
	@interface external:ES~Iterable
	@see {@link
		http://www.ecma-international.org/ecma-262/7.0/index.html#sec-iterable-interface|
		ECMAScript® 2016 Language Specification
		25.1.1.1 The Iterable Interface}
*/

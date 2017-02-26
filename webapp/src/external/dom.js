/**
	@file dom.js includes documentations for DOM.
	This includes citations and you may refer to links for the source.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/**
	DOM defines a platform-neutral model for events and node trees.
	@external DOM
	@see {@link https://dom.spec.whatwg.org/|DOM Standard}
*/

/**
	Interface Event.
	@interface external:DOM~Event
	@extends Object
	@see {@link https://dom.spec.whatwg.org/#event|
		DOM Standard 2.2. Interface Event}
*/

/**
	Interface Element.
	@interface external:DOM~Element
	@extends Object
	@see {@link https://dom.spec.whatwg.org/#interface-element|
		DOM Standard 3.9. Interface Element}
*/

/**
	Interface HTMLElement.
	@interface external:DOM~HTMLElement
	@extends external:DOM~Element
	@see {@link https://html.spec.whatwg.org/multipage/dom.html#htmlelement|
		HTML Standard 3.2.2 Elements in the DOM}
*/

/**
	Interface HTMLInputElement.
	@interface external:DOM~HTMLInputElement
	@extends external:DOM~HTMLElement
	@see {@link https://html.spec.whatwg.org/multipage/forms.html#htmlinputelement|
		HTML Standard 4.10.5 The input element}
*/

/**
	Interface HTMLFormElement.
	@interface external:DOM~HTMLFormElement
	@extends external:DOM~HTMLElement
	@see {@link https://html.spec.whatwg.org/multipage/forms.html#htmlformelement|
		HTML Standard 4.10.3 The form element}
*/

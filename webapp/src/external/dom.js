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
	Callback EventListener.
	@callback external:DOM~EventListener
	@see {@link https://dom.spec.whatwg.org/#eventtarget|
		DOM Standard 2.6. Interface EventTarget}
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
	Interface HTMLImageElement.
	@extends external:DOM~HTMLElement
	@see {@link https://html.spec.whatwg.org/multipage/embedded-content.html#htmlimageelement|
		HTML Standard 4.8.3 The img element}
*/

/**
	Interface HTMLFormElement.
	@interface external:DOM~HTMLFormElement
	@extends external:DOM~HTMLElement
	@see {@link https://html.spec.whatwg.org/multipage/forms.html#htmlformelement|
		HTML Standard 4.10.3 The form element}
*/

/**
	Interface HTMLInputElement.
	@interface external:DOM~HTMLInputElement
	@extends external:DOM~HTMLElement
	@see {@link https://html.spec.whatwg.org/multipage/forms.html#htmlinputelement|
		HTML Standard 4.10.5 The input element}
*/

/**
	Interface CanvasRenderingContext2D.
	@interface external:DOM~CanvasRenderingContext2D
	@extends Object
	@see {@link https://html.spec.whatwg.org/multipage/scripting.html#canvasrenderingcontext2d|
		HTML Standard 4.12.5.1 The 2D rendering context}
*/

/**
	Interface CanvasGradient.
	@interface external:DOM~CanvasGradient
	@extends Object
	@see {@link https://html.spec.whatwg.org/multipage/scripting.html#canvasgradient|
		HTML Standard 4.12.5.1 The 2D rendering context}
*/

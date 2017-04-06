/**
	@file modal.js implements an abstraction for bootstrap modal dialogs
	on Mithril.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/modal */

/**
	modalDialog is an element which will contain the content of the modal
	dialog as its child element.
	@private
	@type !external:jQuery~jQuery
*/
const modalDialog = $(".modal-dialog");

/**
	modal is an element which will be a dialog.
	@private
	@type !external:jQuery~jQuery
*/
const modal = modalDialog.parent();

/**
	front is the front of the list of the queued contents. It is null if the
	list is empty.
	@private
	@type ?module:private/modal~PrivateNode
*/
let front = null;

/**
	beingHidden is the node whose content is mounted to the modal dialog,
	being hidden. It is null if the modal dialog is not being hidden.
	@private
	@type ?module:private/modal~PrivateNode
*/
let beingHidden = null;

/**
	mount is a function to mount a component to the modal dialog.
	@private
	@param {!external:Mithril~Component} - The component to be mounted.
	@returns {Undefined}
*/
function mount(component) {
	m.mount(modalDialog[0], component);
}

/**
	hide is a function to hide the modal dialog. The list should be prepared
	for the content to be shown at the next time AFTER calling this function
	and BEFORE this execution context ends.
	@private
	@returns {Undefined}
*/
function hide() {
	if (!beingHidden) {
		beingHidden = front;

		modal.modal("hide");

		replaceOnhidden(event => {
			const hidden = beingHidden;
			beingHidden = null;

			defaultOnhidden();
			redrawOnhidden(hidden, event);

			if (!hidden.previous && !hidden.next && hidden != front && hidden.component.onmodalremove) {
				hidden.component.onmodalremove(hidden.nodeInterface);
			}
		});
	}
}

/**
	show is a function to show the modal dialog. The component to be
	shown should have been mounted and aria-hidden attribute should have
	been "true" before calling this function.
	@private
	@returns {Undefined}
*/
function show() {
	modal.attr("aria-labelledby", front.attrs["aria-labelledby"]);

	for (const event of ["show", "shown", "hide"]) {
		const bs = event + ".bs.modal";
		const mithril = "onmodal" + event;

		if (front.component[mithril]) {
			modal.one(bs, function() {
				front.component[this](front.nodeInterface);
			}.bind(mithril));
		}
	}

	modal.modal(front.attrs);
}

/**
	draw is a function to mount and show.
	@private
	@returns {Undefined}
*/
function draw() {
	mount(front.component);
	show();
}

/**
	redrawOnhidden prepares for and draws the next entry to be shown if
	it exists. It will clean up if the next entry does not exist.
	@private
	@returns {Undefined}
*/
function redrawOnhidden(previous, event) {
	if (previous.component.onmodalhidden) {
		previous.component.onmodalhidden(event);
	}

	if (front) {
		draw();
	} else {
		mount(null);
		modal.attr("aria-hidden", "true");
	}
}

/**
	init is a function to initialize the modal dialog and the list with the
	given entry.
	@private
	@param {!module:private/modal~PrivateNode} value - The entry which will
	be shown and fill the list.
	@returns {Undefined}
*/
function init(value) {
	front = value;

	if (!beingHidden) {
		modal.attr("aria-hidden", "false");
		draw();
	}
}

/**
	replaceOnhidden is a function to replace the function to be called back
	when "hidden.bs.modal" event gets fired.
	@private
	@param {!external:BS~Handle} callback - The function to
	be called back when "hidden.bs.modal" event gets fired.
	@returns {Undefined}
*/
function replaceOnhidden(callback) {
	modal.off("hidden.bs.modal").on("hidden.bs.modal", callback);
}

/**
	defaultOnhidden is a function to replace the function to be called back
	when "hidden.bs.modal" event with the default function, which will
	show the next entry.
	@private
	@returns {Undefined}
*/
function defaultOnhidden() {
	replaceOnhidden(event => {
		const previous = front;

		front = front.next;
		previous.next = null;

		if (front) {
			front.previous = null;
		}

		redrawOnhidden(previous, event);

		if (previous.component.onmodalremove) {
			previous.component.onmodalremove(previous.nodeInterface);
		}
	});
}

/**
	initNode creates nodeInterface property of the private node and
	calls onmodalinit hook of the component.
	@private
	@param {!module:private/modal~PrivateNode} node - A private node.
	@returns {Undefined}
*/
function initNode(node) {
	node.nodeInterface = {
		remove() {
			if (node == front) {
				hide();
				front = node.next;
			} else if (node.previous) {
				node.previous.next = node.next;
			} else {
				return;
			}

			node.previous = null;
			node.next = null;

			if (beingHidden != node && node.component.onmodalremove) {
				node.component.onmodalremove(node.nodeInterface);
			}
		},

		removed() {
			return !node.previous && node != front;
		},
	};

	if (node.component.onmodalinit) {
		node.component.onmodalinit(node.nodeInterface);
	}
}

defaultOnhidden();

/**
	isEmpty returns a boolean indicating whether the list is empty or not.
	@function
	@returns {!Boolean} a boolean indicating whether the list is empty or
	not.
*/
export const isEmpty = () => !front;

/**
	add adds an entry to the front of the list.
	@function
	@param {!module:private/modal~Attrs} [attrs] - A set of the attributes
	for the new entry.
	@param {!module:private/modal~Component} component - A component for
	the new entry.
	@returns {module:private/modal~Node} A new node.
*/
export function add() {
	const node = {previous: null, next: front};

	if (typeof arguments[0].view == "function") {
		node.attrs = {};
		node.component = arguments[0];
	} else {
		node.attrs = arguments[0];
		node.component = arguments[1];
	}

	if (front) {
		hide();

		front.previous = node;
		front = node;
		initNode(node);
	} else {
		initNode(node);
		init(node);
	}

	return node.nodeInterface;
}

/**
	PrivateNode is a private node in the list.
	@private
	@typedef {Object} module:private/modal~PrivateNode
	@property {!module:private/modal~Attrs} attrs - The attributes of the
	content specified using an exposed API.
	@property {!module:private/modal~Component} component - The content
	which will be mounted to the modal dialog.
	@property {?module:private/modal~PrivateNode} previous - The previous
	node. It is null if the node is removed from the list or it is front.
	@property {?module:private/modal~PrivateNode} next - The next node. It
	is null if the node is removed from the list or it is back.
	@property {?module:private/modal~Node} nodeInterface - The exposed
	interface of this node.
*/

/**
	Node is an object which represents a node in the list.
	@interface module:private/modal~Node
*/

/**
	remove removes the node from the list.
	@function module:private/modal~Node#remove
	@returns {Undefined}
*/

/**
	removed returns a Boolean indicating whether the node is removed from
	the list or not.
	@function module:private/modal~Node#removed
	@returns {!Boolean} A Boolean indicating whether the node is removed
	from the list or not.
*/

/**
	Attrs is a set of attributes for an entry.
	@typedef {external:BS~ModalOptions} module:private/modal~Attrs
	@property {?String} aria-labelledby - The ID of an element which
	represents the label of the entry.
*/

/**
	Component is a component extended to mount to Bootstrap modal dialog.
	@typedef {external:Mithril~Component} module:private/modal~Component

	@property {?module:private/modal~LifecycleMethod} onmodalinit - A hook
	to be called after the component gets added to the list and before the
	component gets mounted.

	@property {?module:private/modal~LifecycleMethod} onmodalshow - A hook
	to be called before the modal dialog gets shown while the component is
	mounted.

	@property {?module:private/modal~LifecycleMethod} onmodalshown - A hook
	to be called after the modal dialog gets shown while the component is
	mounted.

	@property {?module:private/modal~LifecycleMethod} onmodalhide - A hook
	to be called before the modal dialog gets hidden while the component is
	mounted.

	@property {?module:private/modal~LifecycleMethod} onmodalhidden - A hook
	to be called after the modal dialog gets hidden while the component is
	mounted.

	@property {?module:private/modal~LifecycleMethod} onmodalremove - A hook
	to be called after the component gets removed from the list.
*/

/**
	LifecycleMethod is a function to be called in a particular lifecycle
	in the component.
	@callback module:private/modal~LifecycleMethod
	@param {!module:private/modal~Node} node - The node in the list which
	has the component.
*/

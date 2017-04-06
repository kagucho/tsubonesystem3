/**
	@file progress.js implements the control of the progressbar.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/progress */

const progressbar = document.querySelectorAll("[role=progressbar]")[0];

/**
	back is the back of the list of the queued contents. It is null if the
	list is empty.
	@private
	@type ?module:private/progress~PrivateNode
*/
let back = null;

/**
	front is the front of the list of the queued contents. It is null if the
	list is empty.
	@private
	@type ?module:private/progress~PrivateNode
*/
let front = null;

/**
	updateARIA is a function to update the ARIA attributes of the shown
	progress.
	@private
	@returns {Undefined}
*/
function updateARIA() {
	progressbar.setAttribute("aria-describedby", front["aria-describedby"]);
}

/**
	updateValue is a function to update values of the shown progress.
	@private
	@returns {Undefined}
*/
function updateValue() {
	progressbar.setAttribute("aria-valuemax", front.max);
	progressbar.setAttribute("aria-valuenow", front.value);

	if (front.value == front.max) {
		if (!front.important && front.next) {
			front = front.next;
			front.previous.next = null;
			front.previous = null;

			draw();

			return;
		}

		progressbar.setAttribute("aria-busy", "false");
		progressbar.style.animation = ".1s ease 1s 1 normal forwards running progress-hiding";
	}

	progressbar.children[0].style.transform = `scaleX(${front.value / front.max})`;
}

/**
	draw is a function to show the progress at the front of the list.
	@private
	@returns {Undefined}
*/
function draw() {
	progressbar.setAttribute("aria-busy", "true");
	progressbar.style.animationName = "none";
	updateARIA();
	updateValue();
}

/**
	newNode returns a new node.
	@private
	@param {!module:private/progress~Attrs} attrs - A set of the attributes
	of the node.
	@param {!Boolean} important - A boolean indicating whether it is
	important.
	@param {?module:private/progress~Node} previous - The previous node.
	@param {?module:private/progress~Node} next - The next node.
	@returns {!module:private/progress~Node} A new node.
*/
function newNode(attrs, important, previous, next) {
	const node = {
		"aria-describedby": attrs["aria-describedby"],
		max:                attrs.max || 1,
		value:              attrs.value,
		important, previous, next,
	};

	node.nodeInterface = {
		remove() {
			if (node == front) {
				front = front.next;

				if (front) {
					draw();
				} else {
					progressbar.setAttribute("aria-hidden", "true");
					progressbar.style.animation = ".1s ease 1s 1 normal forwards running progress-hiding";
				}
			} else if (node.previous) {
				node.previous.next = node.next;
			} else {
				return;
			}

			if (!node.next) {
				back = node.previous;
			}

			node.previous = null;
			node.next = null;
		},

		removed() {
			return !node.previous && node != front;
		},

		updateARIA(newAttrs) {
			node["aria-describedby"] = newAttrs["aria-describedby"];

			if (node == front) {
				updateARIA();
			}
		},

		updateValue(newAttrs) {
			node.value = newAttrs.value;
			node.max = newAttrs.max || 1;

			if (node == front) {
				updateValue();
			}
		},
	};

	return node;
}

/**
	isEmpty returns a boolean indicating whether the list is empty or not.
	@function
	@returns {!Boolean} a boolean indicating whether the list is empty or
	not.
*/
export const isEmpty = () => !front;

/**
	add is a function to add a new entry.
	@param {!module:private/progress~Attrs} attrs
	@param {?Boolean} important - A boolean indicating the entry is
	important. If it is important, its progress will be kept shown until it
	gets removed. Otherwise, the progress could become hidden even before
	removing.
	@returns {!module:private/progress~Node} A new node in the list of
	the entries.
*/
export function add(attrs, important = false) {
	let node;

	if (front) {
		if (important) {
			if (front.important) {
				node = newNode(attrs, important, front, front.next);
				front.next = node;

				if (node.next) {
					node.next.previous = node;
				}
			} else {
				node = newNode(attrs, important, null, front);
				front.previous = node;
				front = node;

				draw();
			}
		} else {
			node = newNode(attrs, important, back, null);
			back.next = node;
			back = node;
		}
	} else {
		node = newNode(attrs, important, null, null);
		front = node;
		back = node;

		draw();
		progressbar.setAttribute("aria-valuemin", "0");
		progressbar.setAttribute("aria-hidden", "false");
	}

	return node.nodeInterface;
}

/**
	PrivateNode is a private node in the list.
	@private
	@typedef {!module:private/progress~Attrs} module:private/progress~PrivateNode
	@property {!Number} max - The maximum value.
	@property {?module:private/progress~PrivateNode} previous - The
	previous node. It is null if the node is removed from the list or it is
	front.
	@property {?module:private/progress~PrivateNode} next - The next
	node. It is null if the node is removed from the list or it is back.
	@property {?module:private/progress~Node} nodeInterface - The exposed
	interface of this node.
*/

/**
	Node is an object which represents a node in the list.
	@interface module:private/progress~Node
*/

/**
	remove is a function to remove the node from the list.
	@function module:private/progress~Node#remove
	@returns {Undefined}
*/

/**
	removed returns a Boolean indicating whether the node is removed from
	the list or not.
	@function module:private/progress~Node#removed
	@returns {!Boolean} A Boolean indicating whether the node is removed
	from the list or not.
*/

/**
	updateARIA is a function to update the ARIA attributes of the node.
	@function module:private/progress~Node#updateARIA
	@param {!module:private/progress~ARIAAttrs} newAttrs - A set of the
	new ARIA attributes.
	@returns {Undefined}
*/

/**
	updateValue is a function to update the attributes describing values of
	the node.
	@function module:private/progress~Node#updateValue
	@param {!module:private/progress~ValueAttrs} newAttrs - A set of the
	new attributes describing values.
	@returns {Undefined}
*/

/**
	Attrs is an object containing the attributes for an entry.
	@typedef {Object} module:private/progress~Attrs
	@property {!String} aria-labelledby - The ID of an element which
	represents the label of the entry.
	@property {?Number} max - The maximum value. If it is null, it will be
	considered as 1.
	@property {!Number} value - The current value.
*/

/**
	ARIAAttrs is an object containing ARIA attributes for an entry.
	@typedef {Object} module:private/progress~ARIAAttrs
	@property {!String} aria-labelledby - The ID of an element which
	represents the label of the entry.
*/

/**
	ValueAttrs is an object containing attributes describing values for an
	entry.
	@typedef {Object} module:private/progress~ValueAttrs
	@property {?Number} max - The maximum value. If it is null, it will be
	considered as 1.
	@property {!Number} value - The current value.
*/

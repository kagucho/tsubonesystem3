/**
	@file index.js implements the entry point for the public page.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

export default class {
	constructor() {
		this.loaded = 0;
		this.total = 0;
	}

	add(promise) {
		const context = {sum: this, previous: {loaded: 0, total: 0}};

		promise.progress((function(event) {
			this.sum.loaded += event.loaded - this.previous.loaded;
			this.sum.total += event.total - this.previous.total;

			this.previous = event;
		}).bind(context)).always((function() {
			this.sum.loaded -= this.previous.loaded;
			this.sum.total -= this.previous.total;
		}).bind(context));
	}

	html() {
		return {max: this.total, value: this.loaded};
	}
}

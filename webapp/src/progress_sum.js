/**
	@file index.js implements the entry point for the public page.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

export default class {
	constructor() {
		this.computable = false;
		this.loaded = 0;
		this.total = 0;
	}

	add(promise) {
		let previous = {loaded: 0, total: 0};

		promise.progress(event => {
			this.computable = true;
			this.loaded += event.loaded - previous.loaded;
			this.total += event.total - previous.total;

			previous = event;
		}).always(() => {
			this.loaded -= previous.loaded;
			this.total -= previous.total;
		});
	}

	html() {
		return this.total == this.loaded ?
			{value: this.computable ? 1 : 0} :
			{max: this.total, value: this.loaded};
	}
}

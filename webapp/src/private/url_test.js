/**
	@file url_test.js implements a testing code for module:url
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module private/url_test */

import * as url from "./url.js";

describe("module:url", () => {
	describe("mailto", () => {
		it("should return the URL of the given address", () => {
			const expected = "mailto:%20%3F";
			const result = url.mailto(" ?");

			if (result != expected) {
				throw new Error("expected " + expected + ", got " + result);
			}
		});
	});

	describe("tel", () => {
		it("should return the URL of the given telphone number", () => {
			const expected = "tel:+81-3-3260-4271";
			const result = url.tel("03-3260-4271");

			if (result != expected) {
				throw new Error("expected " + expected + ", got " + result);
			}
		});
	});
});

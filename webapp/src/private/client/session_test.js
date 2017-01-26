/**
	@file session_test.js implements a testing code for session.js.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
 */

/** @module private/client/session_test */

import Session from "./session";

describe("module:session", () => {
	const session = new Session;
	const signinPromise = session.signin("1stDisplayID", "1stPassword");
	const sessionPromise = signinPromise.then(() => {
		const result = sessionStorage.getItem("refresh_token");

		if (!result) {
			throw new Error("got "+JSON.stringify(result));
		}
	});

	describe("applyToken", () => {
		const tokenPromise = signinPromise.then(() => session.applyToken(
			token => $.Deferred().resolve(token).promise()));

		it("should return a successful promise", function(done) {
			signinPromise.then(
				() => tokenPromise.then(() => done(), done),
				this.skip.bind(this));
		});

		it("should give callback access_token", function(done) {
			tokenPromise.then(
				token => done(token ? null : Error("got "+JSON.stringify(token))),
				this.skip.bind(this));
		});
	});

	describe("constructor", () => it("should freeze", () => {
		const result = Object.isFrozen(session);
		if (!result) {
			throw new Error("Object.isFrozen(session) returns "+JSON.stringify(result));
		}
	}));

	describe("recover", () => {
		const recovering = new Session;
		const recoverPromise = sessionPromise.then(() => recovering.recover());

		it("should return a successful promise", function(done) {
			sessionPromise.then(() => recoverPromise.done(done), () => this.skip());
		});

		it("should allow to get access_token with applyToken", function(done) {
			recoverPromise.then(() => recovering.applyToken(token =>
				$.Deferred().resolve(token).promise()
			)).then(
				token => done(token ? null : new Error("got "+JSON.stringify(token))),
				this.skip.bind(this));
		});

		it("should allow to get ID with getID", function(done) {
			recoverPromise.then(() => {
				const expected = "1stDisplayID";
				const result = recovering.getID();

				done(result != expected && new Error("expected \"" + expected + "\", got \"" + result + "\""));
			}, this.skip.bind(this));
		});
	});

	describe("signin", () => {
		it("should return a successful promise",
			 done => void signinPromise.then(() => done(), done));

		it("should set refresh_token to sessionStorage", function(done) {
			signinPromise.then(
				() => sessionPromise.then(() => done(), done),
				this.skip.bind(this));
		});

		it("should allow to get ID with getID", function(done) {
			signinPromise.then(() => {
				const expected = "1stDisplayID";
				const result = session.getID();

				done(result != expected && new Error("expected \"" + expected + "\", got \"" + result + "\""));
			}, this.skip.bind(this));
		});
	});
});

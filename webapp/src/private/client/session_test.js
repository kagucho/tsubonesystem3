/**
 * @file session_test.js implements a testing code for session.js.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/** @module session_test */

import Session from "./session.js";

describe("Session", () => {
  const session = new Session;
  const signinPromise = session.signin("1stDisplayId", "1stPassword");
  const sessionPromise = signinPromise.then(() => {
    const result = sessionStorage.getItem("refresh_token");

    if (!result)
      throw new Error("got " + JSON.stringify(result));
  });

  describe("applyToken", () => {
    const tokenPromise = signinPromise.then(() => session.applyToken(token =>
      $.Deferred().resolve(token).promise()));

    it("should return a successful promise", function(done) {
      signinPromise.then(() => tokenPromise.then(() => done(), done),
                         () => this.skip());
    });

    it("should give callback access_token", function(done) {
      tokenPromise.then(
        token => done(token ? null : Error("got " + JSON.stringify(token))),
        () => this.skip());
    });
  });

  describe("constructor", () => it("should freeze", () => {
    const result = Object.isFrozen(session);
    if (!result)
      throw new Error("Object.isFrozen(session) returns " +
                      JSON.stringify(result));
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
        token => done(token ? null : new Error("got " + JSON.stringify(token))),
        () => this.skip());
    });
  });

  describe("signin", () => {
    it("should return a successful promise",
       done => void signinPromise.then(() => done(), done));

    it("should set refresh_token to sessionStorage", function(done) {
      signinPromise.then(() => sessionPromise.done(done), () => this.skip());
    });
  });
});

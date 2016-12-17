/**
 * @file index.js is the main file of client module.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/** @module client */

import * as api from "./api.js";
import Session from "./session.js";

/**
 * state holds private state of client module.
 * @type !external:ES.WeakMap
 */
const state = new WeakMap();

/** module:client is a class to implement the client functionality.
 * @extends external:ES.Object
 */
export default class {
  /** constructor constructs a new instance. */
  constructor() {
    state.set(this, new Session());
    Object.freeze(this);
  }

  /**
   * clubDetail returns the details of the club identified with the given ID.
   * @param {!external:ES.String} id - The ID.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
   * details.
   */
  clubDetail(id) {
    return state.get(this).applyToken(token => api.clubDetail(token, id));
  }

  /**
   * clubList returns the clubs.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
   * clubs.
   */
  clubList() {
    return state.get(this).applyToken(token => api.clubList(token));
  }

  /**
   * clubListName returns the names of clubs associated with IDs.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
   * clubs.
   */
  clubListName() {
    return state.get(this).applyToken(token => api.clubListName(token));
  }

  /**
   * memberDetail returns the details of the member identified with the given
   * ID.
   * @param {!external:ES.String} id - The ID.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
   * details.
   */
  memberDetail(id) {
    return state.get(this).applyToken(token => api.memberDetail(token, id));
  }

  /**
   * memberList returns the members.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
   * members.
   */
  memberList() {
    return state.get(this).applyToken(token => api.memberList(token));
  }

  /**
   * officerDetail returns the details of the officer identified with the given
   * ID.
   * @param {!external:ES.String} id - The ID.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
   * details.
   */
  officerDetail(id) {
    return state.get(this).applyToken(token => api.officerDetail(token, id));
  }

  /**
   * officerList returns the officers.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
   * officers.
   */
  officerList() {
    return state.get(this).applyToken(token => api.officerList(token));
  }

  /**
   * recoverSession recovers the session from sessionStorage.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved when
   * recovered.
   */
  recoverSession() {
    return state.get(this).recover();
  }

  /**
   * signin signs in.
   * @param {!external:ES.String} id - The ID.
   * @param {!external:ES.String} password - The password.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved when
   * signed in.
   */
  signin(id, password) {
    return state.get(this).signin(id, password);
  }
}

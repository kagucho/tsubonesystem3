/**
 * @file session.js implements session.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/** @module session */

import * as api from "./api.js";

/**
 * state holds private state of client module.
 * @type !external:ES.WeakMap
 */
const state = new WeakMap();

/** module:session is a class to implement session.
 * @extends external:ES.Object
 */
export default class {
  /** constructor constructs a new instance. */
  constructor() {
    state.set(this, new class {
      constructor() {
        this.accessToken = null;
        this.refreshToken = null;
      }

      storeRefreshToken() {
        sessionStorage.setItem("refresh_token", this.refreshToken);
      }
    });

    Object.freeze(this);
  }

  /**
   * applyToken applies token to callback which returns a promise.
   * @param {module:session~tokenConsumer} callback - The callback.
   * @returns {!external:jQuery.$.Deferred#promise} A promise returned by the
   * callback.
   */
  applyToken(callback) {
    const local = state.get(this);

    return callback(local.accessToken).catch(xhr => {
      if (xhr.status != 401)
        throw xhr;

      return api.getTokenWithRefreshToken(local.refreshToken).progress(
      progress => progress / 2).then(data => {
        local.accessToken = data.access_token;

        if (data.refresh_token) {
          local.refreshToken = data.refresh_token;
          local.storeRefreshToken();
        }

        return this.applyToken(callback).progress(progress => 0.5 + progress);
      });
    });
  }

  /**
   * recover recovers session from sessionStorage.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved when
   * recovered session.
   */
  recover() {
    const refreshToken = sessionStorage.getItem("refresh_token");

    return api.getTokenWithRefreshToken(refreshToken).then(data => {
      const local = state.get(this);

      local.accessToken = data.access_token;
      if (data.refresh_token) {
        local.refreshToken = data.refresh_token;
        local.storeRefreshToken();
      } else {
        local.refreshToken = refreshToken;
      }
    });
  }

  /**
   * signin signs in.
   * @returns {!external:jQuery.$.Deferred#promise} A promise resolved when
   * signed in.
   */
  signin(id, password) {
    return api.getTokenWithPassword(id, password).then(data => {
      const local = state.get(this);

      local.accessToken = data.access_token;
      local.refreshToken = data.refresh_token;
      local.storeRefreshToken();
    });
  }
}

/**
 * tokenConsumer is a callback for applyToken.
 * @param {!external:ES.String} token - The access token.
 * @returns {!external:jQuery.$.Deferred#promise} A promise which may reject
 * with jqXHR.
 * @callback module:session~tokenConsumer
 */

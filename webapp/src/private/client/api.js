/**
 * @file api.js implements a WebAPI interface.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/** @module api */

/**
 * ajax returns fetched and decoded JSON.
 * @param {!external:ES.String} uri - The URI.
 * @param {!external:ES.String} method - The method.
 * @param {?external:ES.String} token - The token.
 * @param {?external:ES.Object} data - The data.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with
 * fetched JSON.
 */
function ajax(uri, method, token, data) {
  const xhr = new XMLHttpRequest();

  const opts = {
    accepts: {json: "application/json; charset=UTF-8"},
    data, dataType: "json", method,
    xhr() {
      return xhr;
    },
  };

  if (token)
    opts.headers = {Authorization: "Bearer " + token};

  const jqXHR = $.ajax(uri, opts);
  const deferred = $.Deferred();

  xhr.onprogress = function(event) {
    this.notify(event.loaded / event.total);
  }.bind(deferred);

  jqXHR.then(deferred.resolve, deferred.reject);

  return deferred;
}

/**
 * getTokenWithPassword returns a token authorized with the given password.
 * @param {!external:ES.String} username - The ID.
 * @param {!external:ES.String} password - The password.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with a
 * token.
 */
export function getTokenWithPassword(username, password) {
  /* eslint-disable camelcase */
  return ajax("/api/v0/token", "POST", null,
              {grant_type: "password", username, password});
  /* eslint-enable camelcase */
}

/**
 * getTokenWithRefreshToken returns a token authorized with the given refresh
 * toekn.
 * @param {!external:ES.String} token - The refresh token.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with a
 * token.
 */
export function getTokenWithRefreshToken(token) {
  /* eslint-disable camelcase */
  return ajax("/api/v0/token", "POST", null,
              {grant_type: "refresh_token", refresh_token: token});
  /* eslint-enable camelcase */
}

/**
 * clubDetail returns the details of the club identified with the given ID.
 * @param {!external:ES.String} token - The access token.
 * @param {!external:ES.String} id - The ID.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
 * details.
 */
export function clubDetail(token, id) {
  return ajax("/api/v0/club/detail", "GET", token, {id});
}

/**
 * clubList returns the clubs.
 * @param {!external:ES.String} token - The access token.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
 * clubs.
 */
export function clubList(token) {
  return ajax("/api/v0/club/list", "GET", token);
}

/**
 * clubListName returns the names of clubs associated with IDs.
 * @param {!external:ES.String} token - The access token.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
 * club names.
 */
export function clubListName(token) {
  return ajax("/api/v0/club/listname", "GET", token);
}

/**
 * memberDetail returns the details of the member identified with the given ID.
 * @param {!external:ES.String} token - The access token.
 * @param {!external:ES.String} id - The ID.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
 * details.
 */
export function memberDetail(token, id) {
  return ajax("/api/v0/member/detail", "GET", token, {id});
}

/**
 * memberList returns the members.
 * @param {!external:ES.String} token - The access token.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
 * members.
 */
export function memberList(token) {
  return ajax("/api/v0/member/list", "GET", token);
}

/**
 * officerDetail returns the details of the officer identified with the given
 * ID.
 * @param {!external:ES.String} token - The access token.
 * @param {!external:ES.String} id - The ID.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
 * details.
 */
export function officerDetail(token, id) {
  return ajax("/api/v0/officer/detail", "GET", token, {id});
}

/**
 * officerList returns the officers.
 * @param {!external:ES.String} token - The access token.
 * @returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
 * officers.
 */
export function officerList(token) {
  return ajax("/api/v0/officer/list", "GET", token);
}

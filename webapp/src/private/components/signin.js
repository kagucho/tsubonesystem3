/**
	@file signin.js implements the feature to sign in.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/signin */

/**
	module:private/components/signin is a component to provide the feature
	to sign in.
	@name module:private/components/signin
	@type external:Mithril~Component
*/

import * as progress from "./progress";
import client from "../client";

/**
	signin signs in.
	@private
	@param {!external:DOM~HTMLFormElement} form - An element representing
	the credentials of the user.
	@returns {Undefined}
*/
function signin(form) {
	const clientSignin = client.signin(form.id.value, form.password.value);

	this.attrs.onloadstart(clientSignin);

	clientSignin.progress(function(progressEvent) {
		this.progress = {
			max:   progressEvent.total,
			value: progressEvent.loaded,
		};
	}.bind(this.state)).catch(function(xhr) {
		this.error = xhr.responseJSON && xhr.responseJSON.error == "invalid_grant" ?
			"残念！！IDもしくはパスワードが違います。" :
			client.error(xhr);

		delete this.progress;
	}.bind(this.state));

	this.state.progress = {value: 0};
}

export function view(node) {
	const style = {fontSize: "2rem", height: "auto", marginTop: "2rem"};

	return [
		this.progress && m(progress, this.progress),
		m("div", {
			className: "jumbotron",
			style:     {
				margin:  "8rem 0",
				padding: "0",
			},
		}, m("div", {className: "container text-center"},
			m("h1", "TsuboneSystem"),
			m("p", "TsuboneSystemは出欠席管理、メンバー管理、簡易メーリングリスト、非常時連絡先参照を目的に作られたシステムです。"),
			m("p", {className: "hidden-xs"})
		)), m("div", {className: "container"},
			this.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), this.error
			), m("form", {
				"aria-label": "Sign in",
				className:    "text-center",
				role:         "dialog",

				onsubmit(event) {
					signin.call(node, event.target);

					return false;
				},

				style: {maxWidth: "40ch", margin: "0 auto"},
			},
				m("input", {
					autocomplete: "username",
					className:    "form-control",
					inputmode:    "verbatim",
					maxlength:    "64",
					name:         "id",

					oncreate(input) {
						input.dom.focus();
					},

					placeholder: "ID",
					style,
				}), m("input", {
					autocomplete: "current-password",
					className:    "form-control",
					inputmode:    "verbatim",
					maxlength:    "128",
					name:         "password",
					placeholder:  "Password",
					style,
					type:         "password",
				}), m("button", {
					className: "btn btn-lg btn-primary btn-block",
					disabled:  this.progress != null,
					style,
				}, "Sign in")
			)
		),
	];
}

/**
	@file signin.js implements the feature to sign in.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/signin */

/**
	module:private/component/signin is a component to provide the feature
	to sign in.
	@name module:private/component/signin
	@type external:Mithril~Component
*/

import * as progress from "../progress";
import client from "../client";

/**
	signin signs in.
	@private
	@param TODO
	@param {!external:DOM~HTMLFormElement} form - An element representing
	the credentials of the user.
	@returns {Undefined}
*/
function signin(node, form) {
	const clientSignin = client.signin(form.id.value, form.password.value);
	const signinProgress = progress.add({
		"aria-describedby": "component-signin-progress",
		value:              0,
	});

	clientSignin.then(() => {
		signinProgress.remove();
		node.state.progress = null;
		m.redraw();
	}, error => {
		signinProgress.updateARIA({"aria-describedby": "component-signin-error"});

		node.state.error = error == "invalid_grant" ?
			"残念！！IDもしくはパスワードが違います。" :
			client.error(error);

		node.state.progress = null;
		m.redraw();
	}, event => signinProgress.updateValue({
		max:   event.total,
		value: event.loaded,
	}));

	node.state.progress = "サインインしています…";

	if (node.attrs.onloadstart) {
		node.attrs.onloadstart(clientSignin);
	}
}

export function view(node) {
	const style = {fontSize: "2rem", height: "auto", marginTop: "2rem"};

	return [
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
				id:        "component-signin-error",
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

				onsubmit: event => {
					signin(node, event.target);

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
		), m("div", {
			"aria-hidden": (!this.progress).toString(),
			id:            "component-signin-progress",
			style:         {display: "none"},
		}, this.progress),
	];
}

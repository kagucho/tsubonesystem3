/**
	@file signin.js implements the feature to sign in.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module private/components/signin */

/**
	module:private/components/signin is a component to provide the feature
	to sign in.
	@name module:private/components/signin
	@type external:Mithril~Component
*/

import * as client from "../client";
import * as progress from "./progress";

/**
	state holds the private state.
	@type external:ES.WeakMap<external:jQuery.$.Deferred#promise, external:ES.Object>
*/
const state = new WeakMap;

export function controller() {
	const current = {
		deferred: $.Deferred(),
		id:       m.prop(""),
		password: m.prop(""),

		signin() {
			client.signin(this.id(), this.password()).then(() => {
				this.deferred.resolve();
			}, xhr => {
				this.error = xhr.status == 400 ? "残念！！IDもしくはパスワードが違います。" : client.error(xhr);
				delete this.progress;
				m.redraw();
			}, event => {
				this.progress = {
					max:   event.total,
					value: event.loaded,
				};

				m.redraw();
			});

			this.progress = {value: 0};
			m.redraw();
		},
	};

	const promise = current.deferred.promise();
	state.set(promise, current);

	return promise;
}

export function view(promise) {
	const current = state.get(promise);

	const style = {fontSize: "2rem", height: "auto", marginTop: "2rem"};

	return [
		current.progress && m(progress, current.progress),
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
			current.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {ariaHidden: "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), current.error
			),
			m("div", {
				ariaLabel: "Sign in",
				className: "text-center",
				role:      "dialog",
				style:     {maxWidth: "40ch", margin: "0 auto"},
			},
				m("input", {
					className:   "form-control",
					maxlength:   "64",
					oninput:     m.withAttr("value", current.id),
					placeholder: "ID",
					style,
				}), m("input", {
					className:   "form-control",
					maxlength:   "64",
					oninput:     m.withAttr("value", current.password),
					placeholder: "Password",
					style,
					type:        "password",
				}), m("input", {
					className: "btn btn-lg btn-primary btn-block",
					disabled:  current.progress != null,
					onclick:   function() {
						this.signin();
					}.bind(current),
					style,
					type:  "button",
					value: "Sign in",
				})
			)
		),
	];
}

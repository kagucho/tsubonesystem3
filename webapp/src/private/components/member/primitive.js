/**
	@file primitive.js provides the primitive elements of the diplay of a
	member.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/member/primitive */

import * as affiliation from "../../affiliation";
import * as client from "../../client";
import * as password from "../password/modal";
import * as url from "../../url";
import {ensureRedraw} from "../../mithril";
import large from "../../large";

/**
	modal returns a view of a modal dialog.
	@param {!external:ES.Boolean} show - The boolean which indicates whether
	the dialog should be shown.
	@param {!external:jQuery~EventHandler} onhidden - The callback which
	will be called when the modal dialog got hidden.
	@param {?external:ES.Object} options - The options for the containing
	element.
	@param {?external:Mithril~Children} content - The content of the dialog.
*/
const modal = (show, onhidden, options, content) => m("div", $.extend({
	ariaHidden: show.toString(),
	className:  "modal fade",
	config:     (function(element, initialized, context) {
		const jquery = $(element);

		if (!initialized) {
			context.onunload = jquery.modal.bind(jquery, "hide");

			jquery.on("hidden.bs.modal",
				ensureRedraw.bind(undefined, this.onhidden));
		}

		jquery.modal(this.action);
	}).bind({onhidden, action: show ? "show" : "hide"}),
	role:     "dialog",
	tabindex: "-1",
}, options),
	m("div", {className: "modal-dialog", role: "document"},
		m("div", {className: "modal-content"}, content)));

/**
	state holds the internal state.
	@type external:ES.WeakMap<!module:private/components/member/primitive~State>
*/
const state = new WeakMap;

/**
	xhrError returns the message for the error of a XHR.
	@param {!external:jQuery~jqXHR} xhr - A XHR where an error occurred.
	@returns {external:ES.String} The message.
*/
const xhrError = xhr => client.error(xhr) || "どうしようもないエラーが発生しました。";

/**
	State is a class which represents the internal state.
*/
class State {
	/**
		The constructor returns a new State.
		@returns {module:private/components/member/primitive~State} The
		new state.
	*/
	constructor() {
		this.callbacks = {};
		this.member = {current: {}, updated: {}};

		this.modal = {
			messages: {},

			cancel(type) {
				if (type == this.showing) {
					delete this.showing;
					this.callbacks.onhide();
				}
			},

			show(type, message) {
				if (!this.showing) {
					this.callbacks.onshow();
				}

				this.showing = type;
				this.messages[type] = message;
			},
		};

		this.validity = {};
	}

	/**
		promptDeletion prompts the deletion.
		@returns {external:ES~Undefined}
	*/
	promptDeletion() {
		let chief = false;
		for (const club of this.member.current.clubs) {
			if (club.chief) {
				chief = true;
			}
		}

		if (chief || this.member.current.positions.length) {
			this.modal.show("done", {
				body: "役職に就いている局員は削除できません。",
			});
		} else {
			this.modal.show("prompting", {
				body:    "よく狙え。お前は1人の人間を殺すのだ。",
				proceed: {
					label:  "ばーん",
					action: this.submitDeletion.bind(this),
				},
			});
		}
	}

	/**
		submit submits the updated properties of the member.
		@param {!external:ES.Object} callback - TODO
		@returns {external:ES~Undefined}
	*/
	submit(callback) {
		const param = {};

		for (const key in this.member.updated) {
			if (this.member.updated[key] != this.member.current[key]) {
				param[key] = this.member.updated[key];
			}
		}

		let clientSubmit;

		if (this.member.id == null) {
			clientSubmit = client.memberCreate;
		} else {
			let clubsChanged = false;
			for (let index = 0; index < this.clubsUpdated.length; index++) {
				if (this.clubsUpdated[index] != this.clubsCurrent[index]) {
					clubsChanged = true;
				}
			}

			if (clubsChanged) {
				param.clubs = Array.from((function *() {
					for (let index = 0; index < this.clubsUpdated.length; index++) {
						if (this.clubsUpdated[index]) {
							yield this.clubs[index].id;
						}
					}
				})()).join(" ");
			}

			clientSubmit = client.memberUpdate;
		}

		this.callbacks.onloadstart(callback(clientSubmit(param).then(response => {
			this.member.current = $.extend({}, this.member.updated);
			this.clubsCurrent = Array.from(this.clubsUpdated);

			const message = response.error == "mail_failure" ?
				{body: "メールの送信に失敗しました。メールアドレスを確認して下さい。"} :
				{body: "送信しました", success: true};

			this.modal.show("done", message);

			return message.success;
		}, xhr => this.modal.show("done", {body: xhrError(xhr)}))));
	}

	/**
		submitDeletion submits the deletion.
		@returns {external:ES~Undefined}
	*/
	submitDeletion() {
		this.modal.show("inprogress", {body: "始末しています…"});

		this.callbacks.onloadstart(client.memberDelete(this.member.id).then(
		() => this.modal.show("done", {
			body:    "始末しました",
			leave:   this.leaveOnInvalid,
			success: true,
		}), xhr => this.modal.show("done", {body: xhrError(xhr)})));
	}

	/**
		promptOBDeclaration prompts OB declaration.
		@returns {external:ES~Undefined}
	*/
	promptOBDeclaration() {
		this.modal.show("prompting", {
			body:    "一度OB宣言した後に取り消せるように実装されていませんし、面倒なので今後もしません。それでも続行しますか?",
			proceed: {
				label:  "ポチッとな",
				action: this.submitOBDeclaration.bind(this),
			},
		});
	}

	/**
		submitOBDeclaration submits OB declaration.
		@returns {external:ES~Undefined}
	*/
	submitOBDeclaration() {
		this.modal.show("inprogress", {body: "送信しています…"});

		this.callbacks.onloadstart(client.memberDeclareOB().then(() => {
			this.modal.show("done", {
				body:    "送信しました",
				success: true,
			});

			this.member.ob = true;
		}, xhr => this.modal.show("done", {body: xhrError(xhr)})));
	}

	/**
		update updates a property of the member.
		@param {!external:ES.String} key - The key of the property.
		@param {!external:DOM~Event} event - The event whose target
		is an input element which represents.
		@returns {external:ES~Undefined}
	*/
	update(key, event) {
		this.member.updated[key] = event.target.value;
		this.validity[key] = "";
		event.target.checkValidity();
	}

	/**
		updateAttributes updates the control according to the
		attributes.
		@param {!module:private/components/member/primitive~Attributes}
		attributes - The attributes.
		@returns {external:ES~Undefined}
	*/
	updateAttributes(attributes) {
		for (const key of ["onloadstart", "onemptied"]) {
			this.callbacks[key] = attributes[key] || $.noop;
		}

		this.leaveOnInvalid = attributes.leaveOnInvalid;
		this.modal.callbacks = {
			onhide: attributes.onmodalhide || $.noop,
			onshow: attributes.onmodalshow || $.noop,
		};

		this.updateMember(attributes.id);
	}

	/**
		updateClub updates whether the member belongs to the given
		club.
		@param {!external:DOM~Event} event - The event whose target is
		an input element which represents the club and the value to set.
		@returns {external:ES~Undefined}
	*/
	updateClub(event) {
		this.clubsUpdated[event.target.value] = event.target.checked;
	}

	/**
		updateID updates the ID.
		@param {!external:DOM~Event} event - The event whose target is
		an input element where the ID is.
		@returns {external:ES~Undefined}
	*/
	updateID(event) {
		this.member.updated.id = event.target.value;

		if (this.members) {
			for (const member of this.members) {
				if (member.id == this.member.updated.id) {
					event.target.setCustomValidity("そのID、ダブってるよ");

					return this.updateValidity("id", event);
				}
			}
		}

		event.target.setCustomValidity("");
		this.validity.id = "";
		event.target.checkValidity();
	}

	/**
		updateMember updates the member.
		@param {?external:ES.String} id - The ID of the member. The
		form will be prepared for a new member if it is null.
		@returns {external:ES~Undefined}
	*/
	updateMember(id) {
		if (id === this.member.id) {
			return;
		}

		this.member.id = id;

		const management = client.getScope().includes("management");

		let memberDetail;
		if (id == null) {
			if (!management) {
				this.member.id = "?";

				return;
			}

			this.member.current = {
				realname:    "",
				gender:      "",
				mail:        "",
				tel:         "",
				clubs:       [],
				affiliation: "",
				entrance:    new Date().getFullYear(),
			};
			$.extend(this.member.updated, this.member.current);

			this.member.updated.id = "";
			this.member.current.nickname = "新入り";
			this.member.updated.nickname = "";

			this.editable = true;

			memberDetail = $.Deferred().resolve();
		} else {
			const context = {control: this, id: this.member.id};

			memberDetail = client.memberDetail(this.member.id).then((function(member) {
				if (this.control.member.id == this.id) {
					this.control.member.current = member;
					$.extend(this.control.member.updated, member);
				}
			}).bind(context), (function(xhr) {
				if (this.control.member.id == this.id) {
					this.control.callbacks.onemptied();
					throw xhrError(xhr);
				}
			}).bind(context));

			this.callbacks.onloadstart(memberDetail);

			this.editable = this.member.id == client.getID();
			this.deletable = management && !this.editable;
		}

		if (this.editable) {
			if (!this.clubListName) {
				this.clubListName = client.clubListName().then(clubs => {
					this.clubs = clubs;
				}, xhr => {
					throw xhrError(xhr);
				});

				this.callbacks.onloadstart(this.clubListName);

				this.callbacks.onloadstart(client.memberList().then(members => {
					this.members = members;

					affiliation.update(new Set((function *() {
						for (const member of members) {
							yield member.affiliation;
						}
					})()));
				}, xhr => {
					throw xhrError(xhr);
				}));
			}

			$.when(memberDetail, this.clubListName).done((function() {
				if (this.control.member.id == this.id) {
					this.control.clubsCurrent = this.control.clubs.map(club => {
						for (const clubCurrent of this.control.member.current.clubs) {
							if (clubCurrent.id == club.id) {
								return true;
							}
						}

						return false;
					});

					this.control.clubsUpdated = Array.from(this.control.clubsCurrent);
				}
			}).bind({control: this, id: this.member.id}));
		}
	}

	/**
		updateNickname updates the nickname.
		@param {!external:DOM~Event} event - An event whose target is an
		input element which represents the value.
		@returns {external:ES~Undefined}
	*/
	updateNickname(event) {
		this.member.updated.nickname = event.target.value;

		if (this.members) {
			for (const member of this.members) {
				if (member.id != this.member.id && member.nickname == this.member.updated.nickname) {
					event.target.setCustomValidity("そのニックネーム、ダブってるよ");

					return this.updateValidity("nickname", event);
				}
			}

			event.target.setCustomValidity("");
			this.validity.nickname = "";
		}
	}

	/**
		updateValidity updates the validity.
		@param {!external:DOM~Event} event - An event whose target is an
		input element which represents the validity.
		@returns {external:ES~Undefined}
	*/
	updateValidity(key, event) {
		this.validity[key] = event.target.validationMessage;
	}
}

/**
	controlSubmit submits and updates the control accordingly.
	@param {!external:DOM~Event} event - The event.
	@returns {external:ES~Undefined}
	@this module:private/components/member/primitive~Control
*/
function controlSubmit(event) {
	if (this.critical || !event.target.form.checkValidity()) {
		return;
	}

	state.get(this).submit(promise => (this.critical = promise.always(() => {
		this.critical = false;
	})));
}

/**
	controller returns a new control in the MVC architecture.
	@returns {!module:private/components/member/primitive~Control} The control.
*/
export function controller() {
	const controlState = new State;
	const control = Object.defineProperties({}, {
		critical: {writable: true},

		updateAttributes: {
			value: controlState.updateAttributes.bind(controlState),
		},
	});

	state.set(control, controlState);

	return control;
}

/**
	headerView returns the view of the header.
	@param {!module:private/components/member/primitve.controller} control -
	The control.
	@returns {!external:Mithril~Children} The view.
*/
export function headerView(control) {
	const {nickname} = state.get(control).member.current;

	return (nickname == null ? "?" : nickname)+"ちゃんの詳細情報";
}

/**
	bodyView returns the view of the body.
	@param {!module:private/components/member/primitve.controller} control -
	The control.
	@returns {!external:Mithril~Children} The view.
*/
export function bodyView(control) {
	const controlState = state.get(control);
	const headerAttributes = {
		style: {
			fontSize:   "2rem",
			fontWeight: "500",
		},
	};

	let records;
	if (controlState.editable) {
		const nickname = {
			th: m("label", $.extend({
				className: "control-label",
				htmlFor:   "component-member-nickname",
			}, headerAttributes), "ニックネーム"),
			td: m("input", controlState.member.updated.nickname == null ? {
				disabled: true,
			} : {
				className: "form-control",
				id:        "component-member-nickname",
				maxlength: "63",
				onchange:  controlState.updateNickname.bind(controlState),
				oninvalid: controlState.updateValidity.bind(controlState, "nickname"),
				required:  true,
				style:     {display: "inline"},
				value:     controlState.member.updated.nickname,
			}),
			validity: controlState.validity.nickname,
		};

		const mail = {
			th: m("label", $.extend({
				className: "control-label",
				htmlFor:   "component-member-mail",
			}, headerAttributes), "メールアドレス"),
			td: m("input", controlState.member.updated.mail == null ? {
				disabled: true,
			} : {
				className: "form-control",
				id:        "component-member-mail",
				maxlength: "255",
				onchange:  controlState.update.bind(controlState, "mail"),
				oninvalid: controlState.updateValidity.bind(controlState, "mail"),
				required:  true,
				style:     {display: "inline"},
				type:      "email",
				value:     controlState.member.updated.mail,
			}),
			validity: controlState.validity.mail,
		};

		records = controlState.member.id == null ? [
			{
				th: m("label", $.extend({
					className: "control-label",
					htmlFor:   "component-member-id",
				}, headerAttributes), "ID"),
				td: m("input", {
					className: "form-control",
					id:        "component-member-id",
					maxlength: "63",
					onchange:  controlState.updateID.bind(controlState),
					oninvalid: controlState.updateValidity.bind(controlState, "id"),

					/*
						URL Standard
						5.2. application/x-www-form-urlencoded serializing
						https://url.spec.whatwg.org/#urlencoded-serializing

						> 0x2A
						> 0x2D
						> 0x2E
						> 0x30 to 0x39
						> 0x41 to 0x5A
						> 0x5F
						> 0x61 to 0x7A
						>
						> Append a code point whose value is byte to output.

						Accept only those characters.
					*/
					pattern: "[*\\-.\\w]*",

					required: true,
					style:    {display: "inline"},
					title:    "英数字と次の記号 \"*\" \"-\" \".\" \"_\"",
					value:    controlState.member.updated.id,
				}),
				validity: controlState.validity.id,
			},
			nickname,
			mail,
		] : [
			{
				th: m("div", headerAttributes, "ID"),
				td: controlState.member.id,
			},
			{
				th: m("div", headerAttributes, "パスワード"),
				td: m("a", {
					className: "btn btn-primary",
					href:      "#!password",
					onclick:   (function() {
						if (large()) {
							this.show("password");

							return false;
						}
					}).bind(controlState.modal),
				}, "パスワードを変更する"),
			},
			nickname,
			{
				th: m("label", $.extend({
					className: "control-label",
					htmlFor:   "component-member-realname",
				}, headerAttributes), "名前"),
				td: m("input", controlState.member.updated.realname == null ? {
					disabled: true,
				} : {
					className: "form-control",
					id:        "component-member-realname",
					maxlength: "63",
					onchange:  controlState.update.bind(controlState, "realname"),
					oninvalid: controlState.updateValidity.bind(controlState, "realname"),
					required:  true,
					style:     {display: "inline"},
					value:     controlState.member.updated.realname,
				}),
				validity: controlState.validity.realname,
			}, {
				th: m("label", $.extend({
					className: "control-label",
					htmlFor:   "component-member-gender",
				}, headerAttributes), "性別"),
				td: m("input", controlState.member.updated.gender == null ? {
					disabled: true,
				} : {
					className: "form-control",
					id:        "component-member-gender",
					list:      "gender",
					maxlength: "63",
					onchange:  controlState.update.bind(controlState, "gender"),
					oninvalid: controlState.updateValidity.bind(controlState, "gender"),
					style:     {display: "inline"},
					value:     controlState.member.updated.gender,
				}),
				validity: controlState.validity.gender,
			},
			mail,
			{
				// TODO: allow to resend confirming mail if necessary.
				th: m("div", headerAttributes, "メール確認"),
				td: controlState.member.updated.confirmation == null ?
					"?" : (controlState.member.updated.confirmation ? "確認済み" : "未確認"),
			}, {
				th: m("label", $.extend({
					className: "control-label",
					htmlFor:   "component-member-tel",
				}, headerAttributes), "電話番号"),
				td: m("input", controlState.member.updated.tel == null ? {
					disabled: true,
				} : {
					className: "form-control",
					id:        "component-member-tel",
					maxlength: "255",
					onchange:  controlState.update.bind(controlState, "tel"),
					oninvalid: controlState.updateValidity.bind(controlState, "tel"),

					/*
						RFC 3986 - Uniform Resource Identifier (URI): Generic Syntax
						https://tools.ietf.org/html/rfc3986#section-2
						2.2.  Reserved Characters

						Allow characters valid in hier-part.
					*/
					pattern: "[!$&'()*+,\\-.;=~\\w]*",

					required: true,
					style:    {display: "inline"},
					title:    "英数字と次の記号 \"!\" \"$\" \"&\" \"'\" \"(\" \")\" \"*\" \"+\" \",\" \"-\" \".\" \"/\" \";\" \"=\" \"~\"",
					type:     "tel",
					value:    controlState.member.updated.tel,
				}),
				validity: controlState.validity.tel,
			}, {
				th: m("div", headerAttributes, "役職"),
				td: controlState.member.updated.positions == null ? "?" :
					controlState.member.updated.positions.map(position =>
						m("div", m("a", {href: "#!officer?id=" + position.id},
							position.name))),
			}, {
				th: m("div", headerAttributes, "所属部"),
				td: controlState.clubs == null ? "?" : controlState.clubs.map((club, index) => [
					m("label", {
						style: {display: "block"},
					},
						m("input", {
							checked:  controlState.clubsUpdated && controlState.clubsUpdated[index],
							disabled: controlState.clubsUpdated == null,
							onchange: controlState.updateClub.bind(controlState),
							type:     "checkbox",
							value:    index,
						}), " ", club.name
					), " ",
				]),
			}, {
				th: m("label", $.extend({
					className: "control-label",
					htmlFor:   "component-member-affiliation",
				}, headerAttributes), "学科"),
				td: m("input", controlState.member.updated.affiliation == null ? {
					disabled: true,
				} : {
					className: "form-control",
					id:        "component-member-affiliation",
					list:      affiliation.id,
					maxlength: "63",
					onchange:  controlState.update.bind(controlState, "affiliation"),
					oninvalid: controlState.updateValidity.bind(controlState, "affiliation"),
					required:  true,
					style:     {display: "inline"},
					value:     controlState.member.updated.affiliation,
				}),
				validity: controlState.validity.affiliation,
			}, {
				th: m("label", $.extend({
					className: "control-label",
					htmlFor:   "component-member-entrance",
				}, headerAttributes), "入学年度"),
				td: m("input", controlState.member.updated.entrance == null ? {
					disabled: true,
				} : {
					className: "form-control",
					id:        "component-member-entrance",
					max:       "2155",
					min:       "1901",
					onchange:  controlState.update.bind(controlState, "entrance"),
					oninvalid: controlState.updateValidity.bind(controlState, "entrance"),
					required:  true,
					style:     {display: "inline"},
					type:      "number",
					value:     controlState.member.updated.entrance,
				}),
				validity: controlState.validity.entrance,
			}, {
				th: m("div", headerAttributes, "OB宣言"),
				td: controlState.member.updated.ob ? "OB宣言済み" : m("input", {
					className: "btn btn-primary",
					disabled:  controlState.member.updated.ob == null,
					onclick:   controlState.promptOBDeclaration.bind(controlState),
					type:      "button",
					value:     "OB宣言する",
				}),
			},
		];
	} else {
		records = [
			{
				th: m("div", headerAttributes, "ID"),
				td: controlState.member.id == null ?
					"?" : controlState.member.id,
			}, {
				th: m("div", headerAttributes, "ニックネーム"),
				td: controlState.member.updated.nickname == null ?
					"?" : controlState.member.updated.nickname,
			}, {
				th: m("div", headerAttributes, "名前"),
				td: controlState.member.updated.realname == null ?
					"?" : controlState.member.updated.realname,
			}, {
				th: m("div", headerAttributes, "性別"),
				td: controlState.member.updated.gender == null ?
					"?" : controlState.member.updated.gender,
			}, {
				th: m("div", headerAttributes, "メールアドレス"),
				td: controlState.member.updated.mail == null ? "?" :
					m("a", {href: url.mailto(controlState.member.updated.mail)},
						controlState.member.updated.mail),
			}, {
				th: m("div", headerAttributes, "メール確認"),
				td: controlState.member.updated.confirmation == null ?
					"?" : (controlState.member.updated.confirmation ? "確認済み" : "未確認"),
			}, {
				th: m("div", headerAttributes, "電話番号"),
				td: controlState.member.updated.tel == null ? "?" :
					m("a", {href: url.tel(controlState.member.updated.tel)},
						controlState.member.updated.tel),
			}, {
				th: m("div", headerAttributes, "役職"),
				td: controlState.member.updated.positions == null ? "?" :
					controlState.member.updated.positions.map(position =>
						m("div", m("a", {href: "#!officer?id=" + position.id},
							position.name))),
			}, {
				th: m("div", headerAttributes, "所属部"),
				td: controlState.member.updated.clubs == null ? "?" :
					controlState.member.updated.clubs.map(club =>
						m("div",
							m("a", {href: "#!club?id=" + club.id},
								club.name),
							club.chief && " (部長)")),
			}, {
				th: m("div", headerAttributes, "学科"),
				td: controlState.member.updated.affiliation == null ?
					"?" : controlState.member.updated.affiliation,
			}, {
				th: m("div", headerAttributes, "入学年度"),
				td: controlState.member.updated.entrance == null ?
					"?" : controlState.member.updated.entrance,
			}, {
				th: m("div", headerAttributes, "OB宣言"),
				td: controlState.member.updated.ob == null ? "?" :
					(controlState.member.updated.ob ? "OB宣言済み" : "(現役部員)"),
			},
		];
	}

	return m("div", {
		className: "table-responsive",
		style:     {display: "inline-block"},
	},
		m("table", {
			className: "table",
			style:     {textAlign: "left"},
		}, records.map(object => object && m("tr",
			object.validity ? {className: "has-error"} : {},
			m("th", {
				className: "component-member-cell",
				style:     {
					paddingTop: "2rem",
					width:      "16rem",
				},
			}, object.th),
			m("td", {className: "component-member-cell"},
				m("div", {
					className: "component-member-data",
				}, object.td),
				" ",
				m("div", {
					style: {
						display:  "inline-block",
						minWidth: "32ch",
					},
				},
					object.validity && m("div", {
						className: "alert alert-danger",
						role:      "alert",
						style:     {
							display:      "inline-block",
							marginTop:    ".5rem",
							marginBottom: "0",
							padding:      ".5rem",
						},
					},
						m("span", {ariaHidden: "true"},
							m("span", {
								className: "glyphicon glyphicon-exclamation-sign",
							}), " "
						), object.validity
					)
				)
			)
		)))
	);
}

/**
	errorView returns the view of a given error message.
	@param {?external:ES.String} message - The error message.
	@returns {!external:Mithril~Children} The view.
*/
export function errorView(message) {
	return [
		m("span", {ariaHidden: "true"},
			m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
			" "
		), message,
	];
}

/**
	successView returns the view of a given success message.
	@param {?external:ES.String} message - The success message.
	@returns {!external:Mithril~Children} The view.
*/
function successView(message) {
	return [
		m("span", {ariaHidden: "true"},
			m("span", {className: "glyphicon glyphicon-ok"}),
			" "
		), message,
	];
}

/**
	modalView returns the view of the modal dialogs.
	@param {!module:private/components/member/primitve.controller} control -
	The control.
	@returns {!external:Mithril~Children} The view.
*/
export function modalView(control) {
	const controlState = state.get(control);

	const newCanceller =
		type => controlState.modal.cancel.bind(controlState.modal, type);

	const view = [
		controlState.modal.messages.prompting && modal(
			controlState.modal.showing == "prompting",
			newCanceller("prompting"), {
				ariaLabelledby: "component-member-removal-title",
				tabindex:       "-1",
			}, [
				m("div", {className: "modal-header"},
					m("button", {
						ariaLabel:      "閉じる",
						className:      "close",
						type:           "button",
						"data-dismiss": "modal",
					}, m("span", {ariaHidden: "true"}, "×")),
					m("div", {
						className: "lead modal-title",
						id:        "component-member-removal-title",
					}, "確認")
				),
				m("div", {className: "modal-body"},
					controlState.modal.messages.prompting.body),
				m("div", {className: "modal-footer"},
					m("button", {
						className:      "btn btn-default",
						type:           "button",
						"data-dismiss": "modal",
					}, "やっぱやめた"),
					m("button", {
						className: "btn btn-danger btn-pirmary",
						onclick:   controlState.modal.messages.prompting.proceed.action,
					}, controlState.modal.messages.prompting.proceed.label)
				),
			]
		), modal(controlState.modal.showing == "inprogress", null,
			{"data-backdrop": "static"},
			m("div", {className: "modal-body"},
				controlState.modal.messages.inprogress &&
					controlState.modal.messages.inprogress.body
			)
		),
	];

	if (controlState.modal.showing == "password") {
		view.push(m(password, {
			onhidden:    newCanceller("password"),
			onloadstart: controlState.callbacks.onloadstart,
		}));
	}

	if (controlState.modal.messages.done) {
		let onhidden;
		let dismiss;

		if (!controlState.modal.messages.done.leave) {
			onhidden = newCanceller("done");
			dismiss = m("button", {
				className:      "btn btn-default",
				"data-dismiss": "modal",
			}, "閉じる");
		} else if (history.length) {
			onhidden = history.back.bind(history);
			dismiss = m("button", {
				className:      "btn btn-default",
				"data-dismiss": "modal",
			}, "戻る");
		} else {
			onhidden = () => m.route("");
			dismiss = m("a", {
				className: "btn btn-default",
				href:      "#",
			}, "トップページへ");
		}

		view.push(modal(controlState.modal.showing == "done", onhidden, {tabindex: "-1"}, [
			m("div", {className: "modal-body"},
				(controlState.modal.messages.success ? successView : errorView)(
					controlState.modal.messages.done.body)
			), m("div", {className: "modal-footer"}, dismiss),
		]));
	}

	return view;
}

/**
	buttonView returns the view of the control.
	@param {!module:private/components/member/primitve.controller} control -
	The control.
	@returns {!external:Mithril~Children} The view.
*/
export function buttonView(control) {
	const buttons = [];
	const controlState = state.get(control);

	if (controlState.deletable) {
		buttons.push(m("button", {
			className: "btn btn-danger",
			disabled:  controlState.member == null,
			onclick:   controlState.promptDeletion.bind(controlState),
			type:      "button",
		}, "削除"));
	}

	if (controlState.editable) {
		buttons.push(m("button", {
			className: "btn btn-primary",
			disabled:  !controlState.clubs || !controlState.members || controlState.critical,
			onclick:   controlSubmit.bind(control),
			type:      "button",
		}, "送信"));
	}

	return buttons;
}

/**
	An Attributes is an object which contains variables specified as
	attributes.
	@typedef module:private/components/member/primitive~Attributes
	@property {!external:ES.String} id - The ID which identifies the
	member.
	@property {?external:ES.Boolean} leaveOnInvalid - The boolean which
	indicates whether it should leave the current page if the form
	gets invalid. The default behavior is NOT to leave.
	@property {?module:private/components/member/primitive~Onloadstart}
	onloadstart - The function which gets called back when a significant
	and asynchronous loading started.
	@property {?module:private/components/member/primitive~Onemptied}
	onemptied - The function which gets called back when the content is
	missing while it is expected to exist.
	@property {?module:private/components/member/primitive~Onmodalhide}
	onmodalhide - The function which gets called back immediately before the
	modal dialog starts being hidden.
	@property {?module:private/components/member/primitive~Onmodalshow}
	onmodalshow - The function which gets called back immediately before the
	modal dialog starts being shown.
*/

/**
	A Control in the MVC architecture.
	@typedef module:private/components/member/primitive~Control
	@property {!module:private/components/member/primitive~UpdateAttributes}
	updateAttributes - The function to update the control according to the
	attributes.
	@property {?external:ES.Boolean} critical - It is true when the view
	cannot be hidden. This property is NOT configurable.
*/

/**
	An Onloadstart is a function which gets called back when a significant
	and asynchronous loading started.
	@callback module:private/components/member/primitive~Onloadstart
	@param {!external:jQuery.$.Deferred#promise} - TODO
*/

/**
	An Onemptied is a function which gets called back when the content is
	missing while it is expected to exist.
	@callback module:private/components/member/primitive~Onemptied
*/

/**
	An Onmodalhide is a function which gets called back immediately before
	the modal dialog starts being hidden.
	@callback module:private/components/member/primitive~Onmodalhide
*/

/**
	An Onmodalshow is a function which gets called back immediately before
	the modal dialog starts being shown.
	@callback module:private/components/member/primitive~Onmodalshow
*/

/**
	UpdateAttributes updates the context according to the given attributes.
	@callback module:private/components/member/primitive~UpdateAttributes
	@param {!module:private/components/member/primitive~Attributes}
	attributes - The attributes.
	@returns {!external:ES~Undefined}
*/

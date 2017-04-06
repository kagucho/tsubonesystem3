/**
	@file primitive.js provides the primitive elements of the diplay of a
	member.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/member/primitive */

import * as affiliationList from "../../affiliation";
import * as alert from "../alert";
import * as modal from "../../modal";
import * as passwordModal from "../password/modal";
import * as progress from "../../progress";
import * as url from "../../url";
import client from "../../client";
import large from "../../large";

/**
	internals holds the internal states.
	@private
	@type WeakMap.<!module:private/component/member/primitive, !module:private/component/member/primitive~Internal>
*/
const internals = new WeakMap;

/**
	rawProperties are properties whose external representations and internal
	representations are same.
	@private
	@type {!String[]}
*/
const rawProperties = [
	"affiliation", "confirmed", "entrance", "gender",
	"mail", "nickname", "ob", "positions",
	"realname", "tel",
];

/**
	openErrorAlert opens an error dialog. TODO
	@private
	@function
	@param {!module:private/modal~Component} specifiedAlert - A component
	to draw an alert.
	@param {...?external:Mithril~Children} children - An error message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openErrorAlert = (specifiedAlert, ...children) => modal.add(
	specifiedAlert(
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-excalmation-sign"}),
			" "),
		...children));

/**
	openOKAlert opens a dialog showing a successful result. TODO
	@private
	@function
	@param {!module:private/modal~Component} specifiedAlert - A component
	to draw an alert.
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openOKAlert = (specifiedAlert, ...children) => modal.add(
	specifiedAlert(
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-ok"}),
			" "),
		...children));

/**
	openBusyAlert opens a dialog showing progress. TODO
	@private
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openBusyAlert() {
	return modal.add({backdrop: "static"}, alert.busy(...arguments));
}

/**
	openPassword opens a dialog to update the password of the user.
	@private
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
const openPassword = () => modal.add(
	{"aria-labelledby": passwordModal.labelID},
	passwordModal.newComponent());

/**
	openPrompt opens a prompting dialog.
	@private
	@param {!module:private/component/member/primitive~ButtonDescriptor}
	proceeding - A description of proceeding button.
	@param {!module:private/component/member/primitive~ButtonDescriptor}
	cancelling - A description of cancelling button.
	@param {...?external:Mithril~Children} children - A message.
	@returns {!module:private/modal~Node} node - A node of the dialog
	in the list of the modal dialog entries.
*/
function openPrompt(proceeding, cancelling, ...children) {
	let button;

	return modal.add({"aria-labelledby": "component-member-prompt-title"}, {
		onmodalremove: cancelling.action,

		onmodalshown() {
			button.focus();
		},

		view() {
			return m("div", {className: "modal-content"},
				m("div", {className: "modal-header"},
					m("button", {
						"aria-label":   "閉じる",
						"data-dismiss": "modal",
						className:      "close",
						type:           "button",
					}, m("span", {"aria-hidden": "true"}, "×")),
					m("div", {
						className: "lead modal-title",
						id:        "component-member-prompt-title",
					}, "確認")),
				m("div", {className: "modal-body"}, ...children),
				m("div", {className: "modal-footer"},
					m("button", {
						"data-dismiss": "modal",
						className:      "btn btn-default",
						type:           "button",

						oncreate(node) {
							button = node.dom;
						},
					}, cancelling.label),
					m("button", {
						className: "btn btn-danger btn-pirmary",
						onclick:   proceeding.action,
					}, proceeding.label)));
		},
	});
}

/**
	localViewRecordWrapper returns a node wrapping a view of a property.
	@private
	@function
	@param {module:private/component/member/primitive~Property}
		property - The described property.
	@param {module:private/component/member/primitive~Record}
		record - The view of the property.
	@returns {external:Mithril~Children} A node wrapping a view of a
	property.
*/
const localViewRecordWrapper = (property, record) => m("div", {
	className: property.validity ?
		"component-password-record has-error" :
		"component-password-record",
},
	m("dt", {
		style: {
			fontSize:   "2rem",
			fontWeight: "500",
		},
	}, record.dt),
	m("dd", {style: {display: "flex", flexWrap: "wrap"}},
		m("div", {style: {margin: ".5rem"}}, record.dd == null ? "?" : record.dd),
		m("div", {
			"aria-hidden": (!property.validity).toString(),
			style:         {
				margin:     ".5rem",
				minWidth:   "32ch",
				visibility: property.validity ?
					"visible" : "hidden",
			},
		},
			m("div", {
				className: "alert alert-danger",
				role:      "alert",
				style:     {
					display: "inline-block",
					margin:  "0",
					padding: ".5rem",
				},
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "),
				property.validity))));

/**
	localViewRecords returns an iterable of records in the view of local.
	@private
	@param {!module:private/component/member/primitive~Internal}
		internal - The internal state.
	@returns {!Iterable} an iterable of records in the view of local.
*/
function *localViewRecords(internal) {
	const {local} = internal;

	if (!local) {
		return;
	}

	function proxy(key, event) {
		this[key](internal, event.target);
	}

	if (local.id) {
		const {id} = local;

		yield localViewRecordWrapper(id, id.input ? {
			dt: m("label", {
				className: "control-label",
				htmlFor:   id.input,
				style:     {fontWeight: "inherit"},
			}, "ID"),
			dd: m("input", {
				autocomplete: "username",
				className:    "form-control",
				disabled:     id.value == null,
				id:           id.input,
				inputmode:    "verbatim",
				maxlength:    "63",
				oncreate:     internal.prepareFocus.bind(internal),
				oninput:      proxy.bind(id, "update"),
				onchange:     proxy.bind(id, "checkValidity"),
				oninvalid:    proxy.bind(id, "updateValidity"),

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
				value:    id.value || "",
			}),
		} : {dt: "ID", dd: id.value});
	}

	if (local.passwordDialog) {
		const {passwordDialog} = local;

		yield localViewRecordWrapper(passwordDialog, {
			dt: "パスワード",
			dd: m("a", {
				className: "btn btn-primary",
				href:      "#!password",

				onclick() {
					if (large()) {
						openPassword();

						return false;
					}
				},

				role: "button",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-pencil"}),
					" "),
				"パスワードを変更する"),
		});
	}

	if (local.passwordInput) {
		const {passwordInput} = local;

		yield localViewRecordWrapper(passwordInput, {
			dt: m("label", {
				className: "control-label",
				htmlFor:   passwordInput.input,
				style:     {fontWeight: "inherit"},
			}, "パスワード"),
			dd: m("input", {
				autocomplete: "new-password",
				className:    "form-control",
				disabled:     passwordInput.value == null,
				id:           passwordInput.input,
				inputmode:    "verbatim",
				maxlength:    "128",
				oninput:      proxy.bind(passwordInput, "update"),
				onchange:     proxy.bind(passwordInput, "checkValidity"),
				oninvalid:    proxy.bind(passwordInput, "updateValidity"),
				pattern:      "[ -~]*",
				required:     true,
				style:        {display: "inline"},
				title:        "ASCII (英数字と空白, 様々な記号を含む. 詳しくはググれ.)",
				type:         "password",
				value:        passwordInput.value,
			}),
		});
	}

	if (local.passwordVerificationInput) {
		const {passwordVerificationInput} = local;

		yield localViewRecordWrapper(passwordVerificationInput, {
			dt: m("label", {
				className: "control-label",
				htmlFor:   passwordVerificationInput.input,
				style:     {fontWeight: "inherit"},
			}, "パスワード再入力"),
			dd: m("input", {
				autocomplete: "new-password",
				className:    "form-control",
				disabled:     passwordVerificationInput.value == null,
				id:           passwordVerificationInput.input,
				inputmode:    "verbatim",
				maxlength:    "128",
				oninput:      proxy.bind(passwordVerificationInput, "update"),
				onchange:     proxy.bind(passwordVerificationInput, "checkValidity"),
				oninvalid:    proxy.bind(passwordVerificationInput, "updateValidity"),
				pattern:      "[ -~]*",
				required:     true,
				style:        {display: "inline"},
				title:        "ASCII (英数字と空白, 様々な記号を含む. 詳しくはググれ.)",
				type:         "password",
				value:        passwordVerificationInput.value,
			}),
		});
	}

	if (local.nickname) {
		const {nickname} = local;

		yield localViewRecordWrapper(nickname, nickname.input ? {
			dt: m("label", {
				className: "control-label",
				htmlFor:   nickname.input,
				style:     {fontWeight: "inherit"},
			}, "ニックネーム"),
			dd: m("input", {
				autocomplete: "nickname",
				className:    "form-control",
				disabled:     nickname.value == null,
				id:           nickname.input,
				maxlength:    "63",
				oninput:      proxy.bind(nickname, "update"),
				onchange:     proxy.bind(nickname, "checkValidity"),
				oninvalid:    proxy.bind(nickname, "updateValidity"),
				pattern:      "[^\"]*",
				required:     true,
				style:        {display: "inline"},
				title:        "\" (ダブルクオート) は使用できません",
				value:        nickname.value || "",
			}),
		} : {dt: "ニックネーム", dd: nickname.value});
	}

	if (local.realname) {
		const {realname} = local;

		yield localViewRecordWrapper(realname, realname.input ? {
			dt: m("label", {
				className: "control-label",
				htmlFor:   realname.input,
				style:     {fontWeight: "inherit"},
			}, "名前"),
			dd: m("input", {
				className: "form-control",
				disabled:  realname.value == null,
				id:        realname.input,
				maxlength: "63",
				oninput:   proxy.bind(realname, "update"),
				onchange:  proxy.bind(realname, "checkValidity"),
				oninvalid: proxy.bind(realname, "updateValidity"),
				required:  true,
				style:     {display: "inline"},
				value:     realname.value || "",
			}),
		} : {dt: "名前", dd: realname.value});
	}

	if (local.mail) {
		const {mail} = local;

		yield localViewRecordWrapper(mail, mail.input ? {
			dt: m("label", {
				className: "control-label",
				htmlFor:   mail.input,
				style:     {fontWeight: "inherit"},
			}, "メールアドレス"),
			dd: m("input", {
				autocomplete: "email",
				className:    "form-control",
				disabled:     mail.value == null,
				id:           mail.input,
				maxlength:    "255",
				oninput:      proxy.bind(mail, "update"),
				onchange:     proxy.bind(mail, "checkValidity"),
				oninvalid:    proxy.bind(mail, "updateValidity"),
				required:     true,
				style:        {display: "inline"},
				type:         "email",
				value:        mail.value || "",
			}),
		} : {
			dt: "メールアドレス",
			dd: mail.value && m("a", {href: url.mailto(mail.value)},
				mail.value),
		});
	}

	if (local.confirmed) {
		const {confirmed} = local;

		yield localViewRecordWrapper(confirmed, {
			dt: "メール確認",
			dd: confirmed.value == null ?
				"?" :
				(confirmed.value ?
					m("div", {
						className: "label label-success",
						style:     {fontSize: "100%"},
					},
						m("span", {"aria-hidden": "true"},
							m("span", {className: "glyphicon glyphicon-ok"}),
							" "),
						"確認済み") :
					(confirmed.input ?
						m("button", {
							className: "btn btn-primary",
							id:        confirmed.input,
							onclick:   internal.confirm.bind(internal),
							type:      "button",
						}, "確認メールを再送信する") :
						m("div", {
							className: "label label-danger",
							style:     {fontSize: "100%"},
						},
							m("span", {"aria-hidden": "true"},
								m("span", {className: "glyphicon glyphicon-remove"}),
								" "),
							"未確認"))),
		});
	}

	if (local.tel) {
		const {tel} = local;

		yield localViewRecordWrapper(tel, tel.input ? {
			dt: m("label", {
				className: "control-label",
				htmlFor:   tel.input,
				style:     {fontWeight: "inherit"},
			}, "電話番号"),
			dd: m("input", {
				autocomplete: "tel",
				className:    "form-control",
				disabled:     tel.value == null,
				id:           tel.input,
				maxlength:    "255",
				oninput:      proxy.bind(tel, "update"),
				onchange:     proxy.bind(tel, "checkValidity"),
				oninvalid:    proxy.bind(tel, "updateValidity"),

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
				value:    tel.value || "",
			}),
		} : {
			dt: "電話番号",
			dd: tel.value && m("a", {href: url.tel(tel.value)},
				tel.value),
		});
	}

	if (local.gender) {
		const {gender} = local;

		yield localViewRecordWrapper(gender, gender.input ? {
			dt: m("label", {
				className: "control-label",
				htmlFor:   gender.input,
				style:     {fontWeight: "inherit"},
			}, "性別"),
			dd: m("input", {
				autocomplete: "sex",
				className:    "form-control",
				disabled:     gender.value == null,
				id:           gender.input,
				list:         "gender",
				maxlength:    "63",
				oninput:      proxy.bind(gender, "update"),
				onchange:     proxy.bind(gender, "checkValidity"),
				oninvalid:    proxy.bind(gender, "updateValidity"),
				required:     true,
				style:        {display: "inline"},
				value:        gender.value || "",
			}),
		} : {
			dt: "性別",
			dd: gender.value,
		});
	}

	if (local.clubs) {
		const {clubs} = local;

		yield localViewRecordWrapper(clubs, {
			dt: "所属部",
			dd: clubs.value && Array.from(clubs.value, clubs.input ?
				([id, value]) => m("label", {
					style: {display: "block"},
					title: value.chief ? "あんた部長でしょ。辞めさせないわ。" : "",
				},
					m("input", {
						checked:  value.belonging,
						disabled: value.chief,
						onchange: proxy.bind(clubs, "update"),
						type:     "checkbox",
						value:    id,
					}),
					" ",
					value.name,
					value.chief ? [
						" ",
						m("span", {className: "label label-default"}, "部長"),
					] : null) :
				([id, value]) => m("div",
					m("a", {href: "#!club?id=" + id},
						value.name),
					value.chief ? [
						" ",
						m("span", {className: "label label-default"}, "部長"),
					] : null)),
		});
	}

	if (local.positions) {
		const {positions} = local;

		yield localViewRecordWrapper(positions, {
			dt: "役職",
			dd: positions.value && Array.from(positions.value, position =>
				m("div", m("a", {href: "#!officer?id=" + position},
					"TODO: show name"))),
		});
	}

	if (local.ob) {
		const {ob} = local;

		yield localViewRecordWrapper(ob, {
			dt: "OB宣言",
			dd: ob.value == null ?
				"?" :
				(ob.value ?
					"OB宣言済み" :
					(ob.input ?
						m("button", {
							className: "btn btn-primary",
							id:        ob.input,
							onclick:   internal.promptOBDeclaration.bind(internal),
							type:      "button",
						},
							m("span", {"aria-hidden": "true"},
								m("span", {className: "glyphicon glyphicon-check"}),
								" "),
							"OB宣言する") :
						"(現役部員)")),
		});
	}

	if (local.entrance) {
		const {entrance} = local;

		yield localViewRecordWrapper(entrance, entrance.input ? {
			dt: m("label", {
				className: "control-label",
				htmlFor:   entrance.input,
				style:     {fontWeight: "inherit"},
			}, "入学年度"),
			dd: m("input", {
				className: "form-control",
				disabled:  entrance.value == null,
				id:        entrance.input,
				max:       "2155",
				min:       "1901",
				oninput:   proxy.bind(entrance, "update"),
				onchange:  proxy.bind(entrance, "checkValidity"),
				oninvalid: proxy.bind(entrance, "updateValidity"),
				required:  true,
				style:     {display: "inline"},
				type:      "number",
				value:     entrance.value || "",
			}),
		} : {dt: "入学年度", dd: entrance.value});
	}

	if (local.affiliation) {
		const {affiliation} = local;

		yield localViewRecordWrapper(affiliation, affiliation.input ? {
			dt: m("label", {
				className: "control-label",
				htmlFor:   affiliation.input,
				style:     {fontWeight: "inherit"},
			}, "学科"),
			dd: m("input", {
				className: "form-control",
				disabled:  affiliation.value == null,
				id:        affiliation.input,
				list:      affiliationList.id,
				maxlength: "63",
				oninput:   proxy.bind(affiliation, "update"),
				onchange:  proxy.bind(affiliation, "checkValidity"),
				oninvalid: proxy.bind(affiliation, "updateValidity"),
				required:  true,
				style:     {display: "inline"},
				value:     affiliation.value || "",
			}),
		} : {dt: "学科", dd: affiliation.value});
	}
}

/**
	Property is a class representing a property of a member.
	@private
	@extends Object
*/
class Property {
	/**
		constructor constructs module:private/component/member/primitive~Property.
		@param {?String} input - the value to be set to
		input property.
		@param {?*} value - The value to be set to value property.
		@returns {Undefined}
	*/
	constructor(input, value) {
		this.input = input;
		this.value = value;
	}

	/**
		update updates the value of the property.
		@param {!module:private/component/member/primitive~Internal}
		internal - The internal state of the component representing a
		member who has this property.
		@param {!external:DOM~HTMLInputElement} target - An input
		element representing the value of this property.
		@returns {Undefined}
	*/
	update(internal, target) {
		this.value = target.value;
	}

	/**
		input is ID attribute of an input element representing the value
		and the validity of this property. If it is null, it is
		considered readonly.
		@member {?String} module:private/component/member/primitive~Property#input
	*/

	/**
		value is the value of this property. If it is null, it is
		considered not available at the moment.
		@member {?*} module:private/component/member/primitive~Property#value
	*/
}

/**
	ValidatableProperty is a class representing a validatable property of
	a member.
	@private
	@extends module:private/component/member/primitive~Property
*/
class ValidatableProperty extends Property {
	update(internal, target) {
		super.update(internal, target);

		if (target.validity) {
			this.dismissValidity(internal, target);
		}
	}

	/**
		checkValidity checks the validity of his property.
		@param {!module:private/component/member/primitive~Internal}
		internal - The internal state of the component representing a
		member who has this property.
		@param {!external:DOM~HTMLInputElement} target - An input
		element representing the value and validity of this property.
		@returns {!Boolean} - A boolean indicating whether the value is
		valid or not.
	*/
	checkValidity(internal, target) {
		if (this.validity) {
			delete this.validity;
			internal.invalids--;
		}

		return target.checkValidity();
	}

	/**
		dismissValidity checks the validity of this property if
		it was invalid and there is a chance that it will be valid.
		@param {!module:private/component/member/primitive~Internal}
		internal - The internal state of the component representing a
		member who has this property.
		@param {!external:DOM~HTMLInputElement} target - An input
		element representing the value and validity of this property.
		@returns {Undefined}
	*/
	dismissValidity(internal, target) {
		if (this.validity) {
			this.checkValidity(internal, target);
		}
	}

	/**
		reportValidity checks the validity of this property and
		reports the result to the user.
		@param {!module:private/component/member/primitive~Internal}
		internal - The internal state of the component representing a
		member who has this property.
		@param {!external:DOM~HTMLInputElement} target - An input
		element representing the value and validity of this
		property.
		@returns {!Boolean} - A boolean indicating whether the value is
		valid or not.
	*/
	reportValidity(internal, target) {
		if (this.validity) {
			delete this.validity;
			internal.invalids--;
		}

		return target.reportValidity();
	}

	/**
		updateValidity updates the validity of this property.
		@param {!module:private/component/member/primitive~Internal}
		internal - The internal state of the component representing a
		member who has this property.
		@param {!external:DOM~HTMLInputElement} target - An input
		element representing the validity of this property.
		@returns {Undefined}
	*/
	updateValidity(internal, target) {
		if (!this.validity) {
			internal.invalids++;
		}

		this.validity = target.validationMessage;
	}
}

/**
	CustomValidatableProperty is a class representing a vaildatable property
	of a member which have a HTML custom validation.
	@private
	@extends module:private/component/member/primitive~ValidatableProperty
*/
class CustomValidatableProperty extends ValidatableProperty {
	checkValidity(internal, target) {
		this.constructor.setCustomValidity(internal, target);

		return super.checkValidity(internal, target);
	}

	reportValidity(internal, target) {
		this.constructor.setCustomValidity(internal, target);

		return super.reportValidity(internal, target);
	}

	/**
		setCustomValidity sets HTML custom validation.
		@function setCustomValidity
		@memberof module:private/component/member/primitive~CustomValidatableProperty
		@static
		@param {!module:private/component/member/primitive~Internal}
		internal - The internal state of the component representing a
		member who has this property.
		@param {!external:DOM~HTMLInputElement} target - An input
		element representing the validity of this property.
		@returns {Undefined}
	*/
}

/**
	Clubs represents the relation between clubs and a member.
	@private
	@extends module:private/component/member/primitive~Property
*/
class Clubs extends Property {
	update(internal, target) {
		this.value.get(target.value).belonging = target.checked;
	}
}

/**
	ID represents ID of a member.
	@private
	@extends module:private/component/member/primitive~CustomValidatableProperty
*/
class ID extends CustomValidatableProperty {
	static setCustomValidity(internal, target) {
		if (internal.members) {
			for (const member of internal.members) {
				if (member.id == target.value) {
					target.setCustomValidity("そのID、ダブってるよ");

					return;
				}
			}
		}

		target.setCustomValidity("");
	}
}

/**
	Mail represents mail of a member.
	@private
	@extends module:private/component/member/primitive~CustomValidatableProperty
*/
class Mail extends CustomValidatableProperty {
	static setCustomValidity(internal, target) {
		let validity;

		try {
			const ascii = punycode.toASCII(target.value);

			// this value depends on the limit of the database.
			validity = ascii.length > 255 ? "メールアドレス長すぎ" : "";
		} catch (exception) {
			validity = "変なメールアドレスだね";
		}

		target.setCustomValidity(validity);
	}
}

/**
	Nickname represents nickname of a member.
	@private
	@extends module:private/component/member/primitive~CustomValidatableProperty
*/
class Nickname extends CustomValidatableProperty {
	static setCustomValidity(internal, target) {
		if (internal.members) {
			for (const member of internal.members) {
				if (member.id != internal.local.id.value && member.nickname == internal.local.nickname.value) {
					target.setCustomValidity("そのニックネーム、ダブってるよ");

					return;
				}
			}

			target.setCustomValidity("");
		}
	}
}

/**
	PasswordInput represents passwordInput of a member.
	@private
	@extends module:private/component/member/primitive~CustomValidatableProperty
*/
class PasswordInput extends ValidatableProperty {
	dismissValidity(internal, target) {
		super.dismissValidity(internal, target);
		internal.local.passwordVerificationInput.dismissValidity(internal,
			target.form["component-member-password-verification"]);
	}
}

/**
	PasswordVerificationInput represents passwordVerificationInput of a member.
	@private
	@extends module:private/component/member/primitive~CustomValidatableProperty
*/
class PasswordVerificationInput extends CustomValidatableProperty {
	static setCustomValidity(internal, target) {
		target.setCustomValidity(internal.local.passwordInput.value == target.value ?
			"" : "違うよ");
	}
}

/**
	Internal is a class which represents the internal state.
	@private
	@extends Object
*/
class Internal {
	constructor(external) {
		this.external = external;
		this.title = "?";
	}

	/**
		confirm sends an email including the token to confirm the mail.
		@returns {Undefined}
	*/
	confirm() {
		const pullingBusy = openBusyAlert("必要な情報を受信しています…");

		this.pulling.then(() => {
			pullingBusy.remove();

			const submissionModal = openBusyAlert("送信しています…");
			const submissionProgress = progress.add({
				"aria-describedby": alert.bodyID,
				value:              0,
			});

			client.patchUser({mail: this.remote.mail}).then(response => {
				submissionModal.remove();

				if (response.error == "mail_failure") {
					openErrorAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
						"メールの送信に失敗しました。メールアドレスが間違っていないかなど確認してください。");
				} else {
					openOKAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
						`メールを${this.remote.mail}に送信しました。12時間以内に確認してください。ほら早く!`);
				}
			}, error => {
				submissionModal.remove();
				openErrorAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
					client.error(error));
			}, event => submissionProgress.updateValue({
				max:   event.total,
				value: event.loaded,
			}));
		}, () => this.closeInprogress.bind(this));
	}

	/**
		TODO
	*/
	emptyAlert() {
		if (this.external.onemptied) {
			let childIndex;
			let nextOnclosed;

			if (arguments[0] && arguments[0].onclosed) {
				childIndex = 1;
				nextOnclosed = arguments[0].onclosed;
			}

			return alert.closable({
				onclosed() {
					if (nextOnclosed) {
						nextOnclosed();
					}

					this.external.onemptied();
				},
			}, Array.prototype.slice.call(arguments, childIndex));
		} else {
			return alert.leavable(...arguments);
		}
	}

	/**
		prepareFocus prepares to focus.
		@param {!external:Mithril~Node} node - The node to be focused.
		@returns {Undefined}
	*/
	prepareFocus(node) {
		if (this.focusState == "pending") {
			if (document.activeElement == document.body) {
				node.dom.focus();
				this.focusState = null;
			}
		} else {
			this.focusState = node.dom;
		}
	}

	/**
		focus queues focus. TODO
		@returns {!Boolean} A boolean indicating focusing was valid or
		not.
	*/
	focus() {
		if (this.focusState) {
			if (this.focusState.focus) {
				this.focusState.focus();
			}

			return true;
		}

		if (this.local.id) {
			this.focusState = "pending";
			return true;
		}

		return false;
	}

	/**
		isChief returns whether the member is a chief of any club.
		@returns {!Boolean} A boolean indicating the member is a chief
		of any club or not.
	*/
	isChief() {
		for (const {chief} of this.local.clubs.value.values()) {
			if (chief) {
				return true;
			}
		}

		return false;
	}

	/**
		promptDeletion prompts the deletion.
		@returns {Undefined}
	*/
	promptDeletion() {
		if (this.isChief() || this.local.positions.value.length) {
			openErrorAlert(alert.closable, "役職に就いている局員は削除できません。");
		} else {
			openPrompt({
				label:  "ばーん",
				action: this.submitDeletion.bind(this),
			}, {
				label:  "やっぱやめる",
				action: null,
			}, "よく狙え。お前は1人の人間を殺すのだ。");
		}
	}

	/**
		promptOBDeclaration prompts OB declaration.
		@returns {Undefined}
	*/
	promptOBDeclaration() {
		openPrompt({
			label:  "ポチッとな",
			action: this.submitOBDeclaration.bind(this),
		}, {
			label:  "やっぱやめる",
			action: null,
		}, "一度OB宣言した後に取り消せるように実装されていませんし、面倒なので今後もしません。それでも続行しますか?");
	}

	/**
		submit submits the updated properties of the member.
		@param {!external:DOM~HTMLFormElement} target - A form
		representing the validity.
		@returns {Undefined}
	*/
	submit(target) {
		for (const key in this.local) {
			const property = this.local[key];

			if (property.reportValidity && property.input && !property.reportValidity(this, target[property.input])) {
				return;
			}
		}

		let clientSubmit;
		let clubs;
		const param = {};

		switch (this.form) {
		case "creation":
			clientSubmit = client.createMember.bind(client, this.local.id.value);

			for (const key in this.local) {
				if (key != "id") {
					param[key] = this.local[key].value;
				}
			}

			break;

		case "edition":
			for (const key of rawProperties) {
				const property = this.local[key];

				if (property && property.value != this.remote[key]) {
					param[key] = property.value;
				}
			}

			{
				let joining;
				clubs = new Set((function *(internal) {
					for (const [key, value] of internal.local.clubs.value) {
						if (value.belonging) {
							if (!joining) {
								joining = !internal.remote.clubs.has(key);
							}

							yield key;
						}
					}
				})(this));

				if (joining || clubs.size != this.remote.clubs.size) {
					param.clubs = Array.from(clubs);
				}
			}

			clientSubmit = client.patchUser;
			break;

		case "fill":
			for (const key of [
				"affiliation", "entrance", "gender",
				"realname", "tel",
			]) {
				param[key] = this.local[key].value;
			}

			param.clubs = Array.from((function *(clubs) {
				for (const [key, value] of clubs) {
					if (value.belonging) {
						yield key;
					}
				}
			})(this.local.clubs.value));
			param.new_password = this.local.passwordInput.value;

			clientSubmit = client.patchUser;
			break;

		default:
			throw new Error(`unknown form: "${this.form}"`);
		}

		const submissionModal = openBusyAlert("送信しています…");
		const submissionProgress = progress.add({
			"aria-describedby": alert.bodyID,
			value:              0,
		});

		clientSubmit(param).then(response => {
			submissionModal.remove();

			switch (this.form) {
			case "creation":
				if (response.error == "mail_failure") {
					openErrorAlert(this.emptyAlert.bind(this, {onclosed: submissionProgress.remove}),
						"メールの送信に失敗しました。メールアドレスを確認してください。");
				} else {
					openOKAlert(this.emptyAlert.bind(this, {onclosed: submissionProgress.remove}),
						`メールを${param.mail}に送信しました。12時間以内に確認してください。ほら早く!`);
				}
				break;

			case "fill":
				openOKAlert(this.emptyAlert.bind(this, {onclosed: submissionProgress.remove}),
					"神楽坂一丁目通信局へようこそ!");
				break;

			case "edition":
				if (response.error == "mail_failure") {
					openErrorAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
						"メールの送信に失敗しました。メールアドレスを確認してください。");
				} else {
					openOKAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
						param.mail ?
							`メールを${param.mail}に送信しました。12時間以内に確認してください。ほら早く!` :
							"送信しました。");
				}
				break;

			default:
				throw new Error(`unknown form: "${this.form}"`);
			}
		}, error => {
			submissionModal.remove();
			openErrorAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
				client.error(error));
		}, event => submissionProgress.updateValue(
			{max: event.total, value: event.loaded}));
	}

	/**
		submitDeletion submits the deletion.
		@returns {Undefined}
	*/
	submitDeletion() {
		const submissionModal = openBusyAlert("始末しています…");
		const submissionProgress = progress.add({
			"aria-describedby": alert.bodyID,
			value:              0,
		});

		client.deleteMember(this.local.id.value).then(() => {
			submissionModal.remove();
			openOKAlert(this.emptyAlert.bind(this, {onclosed: submissionProgress.remove}),
				"始末しました。");
		}, error => {
			submissionModal.remove();
			openErrorAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
				client.error(error));
		}, event => submissionProgress.updateValue(
			{max: event.total, value: event.loaded}));
	}

	/**
		submitOBDeclaration submits OB declaration.
		@returns {Undefined}
	*/
	submitOBDeclaration() {
		const submissionModal = openBusyAlert("送信しています…");
		const submissionProgress = progress.add({
			"aria-describedby": alert.bodyID,
			value:              0,
		});

		client.declareOB().then(() => {
			submissionModal.remove();
			openOKAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
				"送信しました。");
		}, error => {
			submissionModal.remove();
			openErrorAlert(alert.closable.bind(alert, {onclosed: submissionProgress.remove}),
				client.error(error));
		}, event => submissionProgress.updateValue({
			max:   event.total,
			value: event.loaded,
		}));
	}

	/**
		load loads the details of the member. TODO
		@param {?String} id - The ID of the member. The
		form will be prepared for a new member if it is null.
		@returns {Undefined}
	*/
	start(id) {
		this.invalids = 0;
		this.streams = [];

		const management = client.getScope().includes("management");

		if (id == null) {
			if (!management) {
				throw new Error("TODO: deal with a case that users other than managers tried to create a new member");
			}

			this.local = {
				id:       new ID("component-member-id", ""),
				mail:     new Mail("component-member-mail", ""),
				nickname: new Nickname("component-member-nickname", ""),
			};
			this.form = "creation";
			this.title = "新入り";
		} else {
			const isUser = id == client.getID();
			const local = {
				affiliation: new ValidatableProperty(isUser ? "component-member-affiliation" : undefined),
				clubs:       new Clubs(isUser ? "component-member-clubs" : undefined),
				entrance:    new ValidatableProperty(isUser ? "component-member-entrance" : undefined),
				gender:      new ValidatableProperty(isUser ? "component-member-gender" : undefined),
				id:          new ID(undefined, id),
				realname:    new ValidatableProperty(isUser ? "component-member-realname" : undefined),
				tel:         new ValidatableProperty(isUser ? "component-member-tel" : undefined),
			};

			this.deletable = management && !isUser;
			this.local = local;

			if (isUser) {
				if (client.getFilling()) {
					local.affiliation.value = "";
					local.entrance.value = "";
					local.gender.value = "";
					local.realname.value = "";
					local.tel.value = "";

					local.passwordInput = new PasswordInput("component-member-password", "");
					local.passwordVerificationInput = new PasswordVerificationInput("component-member-password-verification", "");

					this.alert = "これは局員向けのWebサービスです。何にもないですけど、とりあえずフォームを埋めて登録を完了させましょう。今後のサークル生活がフォースと共にあらんことを祈っています。- ある \"TsuboneSystem\" 開発者より";
					this.form = "fill";
					this.pulling = client.mapUser(
						promise => promise.done(member => {
							this.remote = {};
							this.title = member.nickname;
							m.redraw();
						}));
				} else {
					local.confirmed = new Property();
					local.mail = new ValidatableProperty("component-member-mail");
					local.nickname = new ValidatableProperty("component-member-nickname");
					local.ob = new Property();
					local.passwordDialog = new Property();
					local.positions = new Property();

					this.form = "edition";
					this.pulling = client.mapUser(
						promise => promise.done(member => {
							const remote = {clubs: new Set};

							for (const key of rawProperties) {
								if (local[key]) {
									local[key].value = member[key];
								}
							}

							if (!local.clubs.value) {
								local.clubs.value = new Map;
							}

							for (const {id: clubID, chief} of member.clubs()) {
								const old = local.clubs.value.get(clubID);
								if (old) {
									old.belonging = true;
									old.chief = chief;
								} else {
									local.clubs.value.set(clubID, {
										belonging: true,
										chief,
									});
								}

								remote.clubs.add(clubID);
							}

							if (!member.confirmed) {
								local.confirmed.input = "component-member-confirmed";
							}

							if (!member.ob) {
								local.ob.input = "component-member-ob";
							}

							this.remote = $.extend(remote, member);
							this.title = member.nickname;
							m.redraw();
						}));
				}
			} else {
				local.confirmed = new Property();
				local.mail = new Property();
				local.ob = new Property();
				local.positions = new Property();

				this.form = null;
				this.pulling = client.mapMember(id,
					promise => promise.done(member => {
						for (const key of rawProperties) {
							if (local[key]) {
								local[key].value = member[key];
							}
						}

						if (!local.clubs.value) {
							local.clubs.value = new Map;
						}

						for (const {id: clubID, chief} of member.clubs()) {
							const old = local.clubs.value.get(clubID);

							if (old) {
								old.belonging = true;
								old.chief = chief;
							} else {
								local.clubs.value.set(clubID, {
									belonging: true,
									chief,
								});
							}
						}

						this.title = member.nickname;
						m.redraw();
					}));
			}

			this.streams.push(this.pulling,
				client.mapClubs(
					promise => promise.done(
						clubs => {
							const newClubs = new Map;

							for (const club of clubs()) {
								newClubs.set(club.id,
									$.extend(local.clubs.value && local.clubs.value.get(club.id), {
										chief: id == club.chief,
										name:  club.name,
									}));
							}

							local.clubs.value = newClubs;
						})));
		}

		if (this.requiresMembersForVerification()) {
			this.streams.push(affiliationList.listen($.noop));
		}

		client.merge(...this.streams).map(promise => {
			this.loading = "component-member-loading";

			const loadingProgress = progress.add({
				"aria-describedby": this.loading,
				value:              0,
			});

			promise.then(() => {
				this.loading = null;
				loadingProgress.remove();
			}, error => {
				this.loading = null;

				loadingProgress.updateARIA(
					{"aria-describedby": alert.bodyID});

				openErrorAlert(
					alert.closable.bind(alert, {onclosed: loadingProgress.remove}),
					client.error(error));
			}, event => loadingProgress.updateValue(
				{max: event.total, value: event.loaded}));
		});

		m.redraw();
	}

	/**
		TODO
	*/
	end() {
		for (const stream of this.streams) {
			stream.end(true);
		}
	}

	/**
		requiresMembersForVerification returns whether it requires
		members to verify.
		@return {!Boolean} A boolean indicating whether it requires
		members to verify or not.
	*/
	requiresMembersForVerification() {
		return (this.local.id && this.local.id.editable) ||
			(this.local.nickname && this.local.nickname.editable);
	}

	headerView() {
		return this.title + "ちゃんの詳細情報";
	}

	bodyView() {
		return m("div", {className: "text-center"},
			this.alert && m("div", {
				className: "alert alert-info text-left",
				role:      "alert",
			}, this.alert),
			m("dl", {className: "dl-horizontal text-left"},
				...localViewRecords(this)),
			m("div", {
				"aria-hidden": (!this.loading).toString(),
				id:            this.loading,
				style:         {display: "none"},
			}, "読み込み中…"));
	}

	buttonView() {
		const buttons = [];

		if (this.deletable) {
			buttons.push(m("button", {
				className: "btn btn-danger",
				disabled:  this.local == null,
				onclick:   this.promptDeletion.bind(this),
				type:      "button",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-erase"}),
					" "),
				"削除"));
		}

		if (this.form) {
			const attrs = {
				className: "btn btn-primary",
				onclick:   event => this.submit(event.target.form),
			};

			if (!this.local || (this.requiresMembersForVerification() && !this.members)) {
				attrs.disabled = true;
				attrs.title = "もう少し待つのだぞ。";
			} else if (this.invalids) {
				attrs.disabled = true;
				attrs.title = "そんな入力内容で大丈夫か?";
			}

			buttons.push(m("button", attrs, "送信"));
		}

		return buttons;
	}

	/**
		alert is the alerting message.
		@member {?String} module:private/component/member/primitive~Internal#alert
	*/

	/**
		form is a string representing the current role as a form.
		@member {?String} module:private/component/member/primitive~Internal#form
	*/

	/**
		loading is a Boolean indcating whether it is loading. TODO
		@member {!Boolean} module:private/component/member/primitive~Internal#loading
	*/

	/**
		local is the information of the member updated locally.
		@member {?Object.<String, module:private/component/member/primitive~Property>}
			module:private/component/member/primitive~Internal#local
	*/

	/**
		remote is the information of the member in the representation
		of the remote.
		@member {?*} module:private/component/member/primitive~Internal#remote
	*/

	/**
		members is a list of the members.
		@member {?*} module:private/component/member/primitive~Internal#members
	*/

	/**
		pulling is the main loading
		@member {?module:private/promise}
			module:private/component/member/primitive~Internal#pulling
	*/

	/**
		TODO
		@member ?{external:Mithril/Stream[]}
			module:private/component/member/primitive~Internal#streams
	*/

	/**
		title is the title representing the member.
		@member {?String} module:private/component/member/primitive~Internal#title
	*/
}

/**
	State is a class to expose the interfaces for the state and keep track
	of the state.
	@extends Object
*/
export class State {
	/**
		constructor constructs module:private/component/member/primitive.State.
		@returns {Undefined}
	*/
	constructor() {
		const internal = new Internal(this);

		this.header = {
			view() {
				return internal.headerView();
			},
		};

		this.body = {
			view() {
				return internal.bodyView();
			},
		};

		this.button = {
			view() {
				return internal.buttonView();
			},
		};

		internals.set(this, internal);
	}

	/**
		TODO
	*/
	start(id) {
		internals.get(this).start(id);
	}

	/**
		TODO
	*/
	end() {
		internals.get(this).end();
	}

	/**
		confirm sends an email including the token to confirm the mail.
		@returns {Undefined}
	*/
	confirm() {
		internals.get(this).confirm();
	}

	/**
		focus focuses.
		@returns {!Boolean} A boolean indicating whether the focusing
		was valid or not.
	*/
	focus() {
		return internals.get(this).focus();
	}

	/**
		getForm returns the type of the form.
		@returns {!String} "creation" if the form is to create a new
		member. "fill" if the form is to fill the initial information of
		the member as a user. "edition" if the form is to revise the
		information.
	*/
	getForm() {
		return internals.get(this).form;
	}

	/**
		header is a component to draw the header.
		@member module:private/component/member/primitive.State#header
	*/

	/**
		body is a component to draw the body.
		@member module:private/component/member/primitive.State#body
	*/

	/**
		button is a component to draw buttons.
		@member module:private/component/member/primitive.State#button
	*/
}

/**
	error is a component providing the view of the error. Let the value
	returned by rejected loading be the children to show the error message.
	@type !external:Mithril~Component
*/
export const error = Object.freeze({
	view(node) {
		return [
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
				" "
			), node.children,
		];
	},
});


/**
	success is a component providing the view of the success message. Let
	the value returned by rejected loading be the children to show the error
	message.
	@type !external:Mithril~Component
*/
export const success = Object.freeze({
	view(node) {
		return [
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-ok"}),
				" "
			), node.children,
		];
	},
});

/**
	ButtonDescriptor describes a button.
	@private
	@typedef module:private/component/member/primitive~ButtonDescriptor
	@property {?function} action - The function to be called after the
	button gets clicked.
	@property {!String} label - The label of the button.
*/

/**
	Record represents a view of a member.
	@private
	@typedef module:private/component/member/primitive~Record
	@property {!external:Mithril~Children} dt - The definition term.
	@property.{?external:Mithril~Children} dd - The definition description.
*/

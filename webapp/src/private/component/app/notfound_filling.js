/**
	@file notfound_filling.js implements notfoundFilling component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/notfoundFilling */

/**
	module:private/component/app/notfoundFilling is a component to show the
	requested page was not found when prompting the user to fill his
	information.
	@name module:private/component/app/notfoundFilling
	@type !external:Mithril~Component
*/

import * as container from "../container";
import client from "../../client";
import logo from "../../../images/logo250c_black.png";

export const view = () => m(container,
	m("img", {
		alt:       "ロゴ",
		className: "hidden-xs pull-right",
		src:       logo,
	}), m("article", {className: "container"},
		m("h1", "未登録者が閲覧できないページです"),
		m("p",
			"まだあなたは情報の入力を完了していませんね?",
			m("a", {href: "#!member?id=" + client.getID()},
				"こちらで入力を完了させてください.")
		)
	)
);

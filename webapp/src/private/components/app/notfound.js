/**
	@file notfound.js implements notfound component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/notfound */

/**
	module:private/components/app/notfound is a component to show the
	requested page was not found.
	@name module:private/components/app/notfound
	@type !external:Mithril~Component
*/

import * as container from "../container";
import logo from "../../../images/logo250c_black.png";

export const view = () => m(container,
	m("img", {
		alt:       "ロゴ",
		className: "hidden-xs pull-right",
		src:       logo,
	}), m("article", {className: "container"},
		m("h1", "ページが見つかりませんでした"),
		m("section",
			m("h2", "どうしよう?"),
			m("dl",
				m("dt", "このサイトにあるリンクを踏んでここに来た。"),
				m("dd",
					m("a", {href: "mailto:kagucho.net@gmail.com"},
						"kagucho.net@gmail.com"
					), "へ報告してください。お願いします。あなたの好意が世界を救います。"
				), m("dt", "自分でURIを入力した。"),
				m("dd", "打ち間違えてないか確認してください。打ち間違えてない? それは困った。お役に立てそうにない。"),
				m("dt", "トライフォースを探しにここに来た。"),
				m("dd", "ここにはないです。本当に。")
			)
		)
	)
);

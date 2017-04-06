/**
	@file mails.js implements mails component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

import * as container from "../container";
import * as modal from "../../modal";
import * as mailModal from "../mail/modal";
import * as progress from "../../progress";
import client from "../../client";
import large from "../../large";
import memberModal from "../member/modal";

/**
	TODO
*/
function openMail(subject) {
	if (large()) {
		this.mailModal = modal.add(
			mailModal.default(subject),
			{"aria-labelledby": mailModal.labelID});

		return false;
	}

	return true;
}

/**
	TODO
*/
function openMember(id) {
	if (large()) {
		this.memberModal = modal.add(memberModal(id));
		return false;
	}

	return true;
}

export function onbeforeremove() {
	if (this.streams) {
		for (const stream of this.streams) {
			stream.end(true);
		}
	}

	if (this.mailModal) {
		this.mailModal.remove();
	}

	if (this.memberModal) {
		this.memberModal.remove();
	}
}

export function oninit() {
	this.streams = [
		client.mapMembers(
			promise => promise.then(
				members => {
					const object = Object.create(null);

					for (const value of members()) {
						object[value.id] = value;
					}

					return object;
				})),

		client.mapMails(
			promise => promise.then(
				mails => {
					this.mails = Array.from(mails());
					m.redraw();
				})),
	];

	client.merge(...this.streams).map(promise => {
		this.loading = "component-app-mails-loading";

		const loadingProgress = progress.add(
			{"aria-describedby": this.loading, value: 0});

		promise.then(members => {
			this.loading = null;

			for (const mail of this.mails) {
				mail.fromInformation = members[value.from];
			}

			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.error = client.error(error);
			this.loading = null;
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	});
}

export function view() {
	return [
		m(container,
			m("div", {className: "container table-responsive"},
				m("h1", "Mails"),
				this.error && m("div", {
					className: "alert alert-danger",
					role:      "alert",
				},
					m("span", {"aria-hidden": "true"},
						m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
						" "),
					this.error),
				m("a", {
					className: "btn btn-block btn-lg btn-primary",
					href:      "#!mail",
					onclick:   openMail.bind(this, null),
					role:      "button",
					style:     {margin: "1rem"},
				}, "メールする!"),
				m("section", {className: "table-responsive"},
					m("h2", "過去のメール"),
					m("p", {
						className: "lead",
						style: {color: "gray"},
					}, this.mails && this.mails.length + "件"),
					m("table", {className: "table"},
						m("thead",
							m("tr",
								m("th", "Date"),
								m("th", "Subject"),
								m("th", "From"),
								m("th", "To"))),
						m("tbody",
							this.mails && this.mails.map(mail => m("tr",
								{key: mail.subject},
								m("td", moment.unix(mail.date).format("lll")),
								m("td",
									m("a", {
										href:    "#!mail?subject=" + encodeURIComponent(mail.subject),
										onclick: openMail.bind(this, mail.subject),
									}, mail.subject)),
								m("td",
									m("a", {
										href:    "#!member?id=" + mail.from,
										onclick: openMember.bind(this, mail.from),
									}, mail.fromInformation && mail.fromInformation.nickname)),
								m("td", mail.to)))))))),
		m("div", {
			"aria-hidden": (!this.loading).toString(),
			id:            this.loading,
			style:         {display: "none"},
		}, "読み込み中…"),
	];
}

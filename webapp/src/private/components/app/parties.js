/**
	@file parties.js implements parties component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

import * as alert from "../alert";
import * as container from "../container";
import * as modal from "../../modal";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";
import client from "../../client";

function openOK() {
	return modal.unshift(alert.closable(
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-ok"}),
			" "
		), ...arguments
	));
}

function openError() {
	return modal.unshift(alert.closable(
		m("span", {"aria-hidden": "true"},
			m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
			" "
		), ...arguments
	));
}

function openInprogress() {
	return modal.unshift({backdrop: "static"},
		alert.inprogress(...arguments));
}

/**
	TODO
*/
function showPlannedParties() {
	this.showingPastParties = false;
}

/**
	TODO
*/
function showPastParties() {
	this.showingPastParties = true;
}

function respond(party, attending) {
	const inprogress = openInprogress("送信しています…");

	this.progress.add(client.partyRespond({
		party:     party.name,
		attending: attending ? "1" : "0",
	}).then(() => {
		party.user = attending ? "accepted" : "declined";

		inprogress.remove();
		openOK("送信しました。");
	}, xhr => {
		inprogress.remove();
		openError(client.error(xhr) || "どうしようもないエラーです。");
	}));
}

/**
	TODO
*/
function partiesView() {
	const now = moment();

	return this.parties.map(party => {
		const start = moment.unix(party.start);
		const end = moment.unix(party.end);
		const due = moment.unix(party.due);

		return this.showingPastParties == now.isAfter(end) && m("tr",
			m("td", party.name),
			m("td", m("dl", {style: {display: "table"}},
				m("div", {style: {display: "table-row"}},
					m("dt", {
						style: {
							display: "table-cell",
							padding: "1rem",
						},
					}, "開始"),
					m("dd", {
						style: {
							display: "table-cell",
							padding: "1rem",
						},
					}, start.format("lll"))
				), m("div", {style: {display: "table-row"}},
					m("dt", {
						style: {
							display: "table-cell",
							padding: "1rem",
						},
					}, "終了"),
					m("dd", {
						style: {
							display: "table-cell",
							padding: "1rem",
						},
					}, end.format("lll"))
				)
			)), m("td", party.place),
			m("td", party.inviteds),
			m("td", due.format("lll")),
			m("td", {
				uninvited: m("span",
					{className: "label label-default"},
					"招待されていません"
				),

				invited: now.isBefore(due) ? [
					m("button", {
						className: "btn btn-block btn-primary",
						onclick:   respond.bind(this, party, true),
					}, "出席"),
					m("button", {
						className: "btn btn-block btn-danger",
						onclick:   respond.bind(this, party, false),
					}, "欠席"),
				] : m("span",
					{className: "label label-warning"},
					"未提出"
				),

				accepted: m("span",
					{className: "label label-primary"},
					"出席予定"
				),

				declined: m("span",
					{className: "label label-danger"},
					"欠席予定"
				),
			}[party.user])
		)
	});
}

export function oninit() {
	const listing = client.partyList().then(
		parties => this.parties = parties,
		xhr => this.error = client.error(xhr) || "どうしようもないエラーです");

	this.error = null;
	this.progress = new ProgressSum;
	this.progress.add(listing);
	this.showingPastParties = m.route.param().past || false;
}

export function view() {
	return [
		m(progress, this.progress.html()),
		m(container, m("div", {className: "container"},
			this.error && m("section", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), this.error
			), m("h1", "Parties"),
			m("a", {
				className: "btn btn-block btn-lg btn-primary",
				href:      "#!party",
				style:     {margin: "1rem"},
			}, "パーティーする!"),
			m("section", {
				className: "panel panel-default",
				style:     {margin: "1rem"},
			},
				m("ul", {className: "nav nav-tabs"},
					m("li", {
						className: this.showingPastParties ? "" : "active",
						role:      "presentation",
					}, m("a", {
						href:    "#!parties",
						onclick: showPlannedParties.bind(this),
					},
						"予定されているパーティー"
					)), m("li", {
						className: this.showingPastParties ? "active" : "",
						role:      "presentation",
					}, m("a", {
						href:    "#!parties?past=1",
						onclick: showPastParties.bind(this),
					},
						"過去のパーティー"
					))
				), m("div", {
					className: "table-responsive",
				}, m("table", {className: "table"},
					m("thead",
						m("tr",
							m("th", "題目"),
							m("th", "時刻"),
							m("th", "場所"),
							m("th", "対象者"),
							m("th", "出欠締め切り"),
							m("th", "出欠")
						)
					), m("tbody", this.parties ?
						partiesView.call(this) : [
							m("span", {"aria-hidden": "true"},
								m("span", {className: "glyphicon glyphicon-hourglass"}),
								" "
							), "読み込み中です…",
						]
					)
				))
			)
		)),
	];
}

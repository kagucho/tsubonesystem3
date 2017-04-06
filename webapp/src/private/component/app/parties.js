/**
	@file parties.js implements parties component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

import * as alert from "../alert";
import * as container from "../container";
import * as modal from "../../modal";
import * as progress from "../../progress";
import client from "../../client";

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

/**
	TODO
*/
function respond(party, attending) {
	const responseBusy =
		modal.add({backdrop: "static"}, alert.busy("送信しています…"));

	const responseProgress = progress.add(
		{"aria-describedby": alert.bodyID, value: 0});

	client.respondParty(party.name, attending).then(() => {
		responseBusy.remove();

		modal.add(alert.closable(
			{onclosed: responseProgress.remove},
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-ok"}),
				" "),
			"送信しました。"));
	}, error => {
		responseBusy.remove();

		modal.add(alert.closable(
			{onclosed: responseProgress.remove},
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
				" "),
			client.error(error)));
	}, event => responseProgress.updateValue(
		{max: event.total, value: event.loaded}));
}

/**
	TODO
*/
function partiesView() {
	const now = moment();

	return Array.from(this.parties(), party => {
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
					}, start.format("lll"))),
				m("div", {style: {display: "table-row"}},
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
				))),
			m("td", party.place),
			m("td", party.inviteds),
			m("td", due.format("lll")),
			m("td", {
				uninvited: m("span",
					{className: "label label-default"},
					"招待されていません"),

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
					"未提出"),

				accepted: m("span",
					{className: "label label-primary"},
					"出席予定"),

				declined: m("span",
					{className: "label label-danger"},
					"欠席予定"),
			}[party.user]));
	});
}

export function oninit() {
	this.endPartiesStream = client.mapParties(promise => {
		const loadingProgress = progress.add({
			"aria-describedby": "component-app-parties-loading",
			value:              0,
		});

		promise.then(parties => {
			this.parties = parties;
			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.error = client.error(error);
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	}).end;

	this.error = null;
	this.showingPastParties = m.route.param().past || false;
}

export function onremove() {
	this.endPartiesStream(true);
}

export function view() {
	return m(container, m("div", {className: "container"},
		this.error && m("section", {
			className: "alert alert-danger",
			role:      "alert",
		},
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
				" "),
			this.error),
		m("h1", "Parties"),
		m("a", {
			className: "btn btn-block btn-lg btn-primary",
			href:      "#!party",
			role:      "button",
			style:     {margin: "1rem"},
		}, "パーティーする!"),
		m("section", {
			className: "panel panel-default",
			style:     {margin: "1rem"},
		},
			m("ul", {
				className: "nav nav-tabs",
				role:      "tablist",
			},
				m("li", {
					className: this.showingPastParties ? "" : "active",
					id:        "component-app-parties-planned",
					role:      "tab",
				}, m("a", {
					href:    "#!parties",
					onclick: showPlannedParties.bind(this),
				}, "予定されているパーティー")),
				m("li", {
					className: this.showingPastParties ? "active" : "",
					id:        "component-app-parties-past",
					role:      "tab",
				}, m("a", {
					href:    "#!parties?past=1",
					onclick: showPastParties.bind(this),
				}, "過去のパーティー"))),
			m("div", {
				"aria-labelledby": this.showingPastParties ?
					"component-app-parties-past" :
					"component-app-parties-planned",
				className: "table-responsive",
				role:      "tabpanel",
			}, m("table", {className: "table"},
				m("thead",
					m("tr",
						m("th", "題目"),
						m("th", "時刻"),
						m("th", "場所"),
						m("th", "対象者"),
						m("th", "出欠締め切り"),
						m("th", "出欠"))),
				m("tbody", this.parties ?
					partiesView.call(this) :
					m("div", {id: "component-app-parties-loading"},
						m("span", {"aria-hidden": "true"},
							m("span", {className: "glyphicon glyphicon-hourglass"}),
							" "),
						"読み込み中です…")))))));
}

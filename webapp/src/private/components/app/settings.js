import * as client from "../../client";
import * as container from "../container";
import * as member from "../member/modal";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";
import large from "../../large";
import {members} from "../table";

export class controller {
	constructor() {
		this.member = {id: client.getID()};
		this.progress = new ProgressSum;
		this.load();

		const confirm = m.route.param("confirm");
		if (confirm) {
			this.progress.add(client.memberConfirm(confirm).then(() => {
				this.confirmation = {
					state:   "done",
					success: true,
					message: "メールアドレスを確認しました。",
				};

				this.load();
			}, xhr => {
				this.confirmation = xhr.responseJSON && xhr.responseJSON.error == "invalid_request" ? {
					state: "abortOrRetry",
				} : {
					state:   "done",
					message: client.error(xhr) || "どうしようもないエラーです。",
				};
			}));

			this.confirmation = {state: "sending"};
		}
	}

	confirm() {
		this.progress.add(client.memberUpdate({mail: this.member.mail}).then(() => {
			this.confirmation = {
				state:   "done",
				success: true,
				message: `メールを${this.member.mail}に送信しました。12時間経過後に無効になるのでさっさと確認してください。`,
			};
		}, xhr => {
			this.confirmation = {
				state:   "done",
				message: client.error(xhr) || "どうしようもないエラーが発生しました。",
			};
		}));

		this.confirmation = {state: "sending"};
	}

	hide() {
		delete this.showing;
	}

	finishConfirmation() {
		delete this.confirmation;
	}

	load() {
		this.progress.add(client.memberDetail(this.member.id).then(
			$.extend.bind($, this.member),
			xhr => {
				this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
			}
		));
	}

	show() {
		if (large()) {
			this.showing = true;

			return false;
		}
	}

	submissionStart(promise) {
		this.progress.add(promise.done(submission => submission && this.load()));
	}
}

export function view(control) {
	return [
		m(progress, control.progress.html()),
		m(container,
			m("div", {className: "container"},
				m("h1", "Settings"),
				m("div", "何して遊ぶ? 何して遊ぶ? んっ? 遊ばないのかぁ"),
				control.error && m("div",
					m("span", {ariaHidden: "true"},
						m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
						" "
					),
					control.error
				),
				control.member && [
					m(members, {
						members:     [control.member],
						onloadstart: control.submissionStart.bind(control),
					}),
					m("a", {
						className: "btn btn-default",
						href:      "#!member?id=" + control.member.id,
						onclick:   control.show.bind(control),
					}, "登録情報の表示, 変更"),
					" ",
					control.member.confirmed ?
						"メールアドレスは確認済みです。やったね!" :
						m("button", {
							className: "btn btn-default",
							onclick:   control.confirm.bind(control),
							type:      "button",
						}, "確認メールの再送信"),
				]
			)
		),
		control.showing && m(member, {
			id:          control.member.id,
			onhidden:    control.hide.bind(control),
			onloadstart: control.submissionStart.bind(control),
		}),
		control.confirmation && [
			{
				state: "done",
				content: m("div", {className: "modal-content"},
					m("div", {className: "modal-body"},
						m("span", {ariaHidden: "true"},
							m("span", {
								className: control.confirmation.success ?
									"glyphicon glyphicon-ok" :
									"glyphicon glyphicon-exclamation-sign",
							}), " "
						), control.confirmation.message
					), m("div", {className: "modal-footer"},
						m("button", {
							className:      "btn btn-default",
							config:         control.confirmation.state == "done" &&
									control.confirmation.shown &&
									(element => element.focus()),
							type:           "button",
							"data-dismiss": "modal",
						}, "閉じる")
					)
				),
			}, {
				state:   "abortOrRetry",
				content: m("div", {className: "modal-content"},
					m("div", {className: "modal-body"},
						"メールアドレスの確認に失敗しました。たぶん時間切れとかそんなところじゃないですかね? 確認メールを再送信しますか?"
					), m("div", {className: "modal-footer"},
						m("button", {
							className:      "btn btn-default",
							type:           "button",
							"data-dismiss": "modal",
						}, "やっぱやめる"),
						m("button", {
							className: "btn btn-primary",

							config: control.confirmation.state == "abortOrRetry" &&
								control.confirmation.shown &&
								(element => element.focus()),

							type:      "button",
							onclick:   control.confirm.bind(control),
						}, "再送信する")
					)
				),
			}, {
				state: "sending",
				content: m("div", {className: "modal-content"},
					m("div", {className: "modal-body"},
						"送信しています…"
					)
				),
			},
		].map(object => m("div", {
			ariaHidden: (control.confirmation.state != object.state).toString(),
			className:  "modal fade",
			config: (function(state, element, initialized, context) {
				const jquery = $(element);

				if (!initialized) {
					context.onunload = jquery.modal.bind(jquery, "hide");

					jquery.on("hidden.bs.modal",
						() => this.confirmation.state == state && this.finishConfirmation());

					jquery.on("shown.bs.modal", () => {
						this.confirmation.shown = true;
						m.redraw();
					});
				}

				jquery.modal(this.confirmation.state == state ? "show" : "hide");
			}).bind(control, object.state),
		}, m("div", {className: "modal-dialog", role: "document"},
			object.content
		))),
	];
}

/**
	@file index_test.js implements a testing code for index.js.
	@author Akihiko Odaki  <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client_test */

import * as client from "./index";

describe("module:client", () => {
	const signinPromise = client.signin("1stDisplayID", "1stPassword");

	describe("signin", () => {
		it("should return a successful promise",
			signinPromise.done.bind(signinPromise));

		it("should set refresh_token to sessionStorage", function(done) {
			signinPromise.then(() => {
				const result = sessionStorage.getItem("refresh_token");
				done(result ? null : new Error("got "+JSON.stringify(result)));
			}, this.skip.bind(this));
		});
	});

	describe("memberUpdate", () => it("should reject with invalid entrance", function(done) {
		signinPromise.then(() => {
			client.memberUpdate({entrance: "0"}).then(
			done.bind(undefined, new Error("unexpected resolution")),
			xhr => {
				const expected = "value out of range";
				done(xhr.responseJSON.error_description != expected &&
					new Error("expected \""+xhr.responseJSON.error_description+
						"\", got \""+expected+"\""));
			});
		}, this.skip.bind(this));
	}));

	describe("recoverSession", () => {
		const recoverPromise = signinPromise.then(() =>
			client.recoverSession());

		it("should return a successful promise", function(done) {
			signinPromise.then(
				() => recoverPromise.done(done),
				this.skip.bind(this));
		});
	});

	const staticTests = [
		{
			name:     "clubDetail",
			run:      () => client.clubDetail("prog"),
			expected: {
				chief: {
					id:       "2ndDisplayId",
					mail:     "",
					nickname: " !%_1\"#",
					realname: "$&\\%_2'(",
					tel:      "",
				},
				members: [
					{
						entrance: 1901,
						id:       "2ndDisplayId",
						nickname: " !%_1\"#",
						realname: "$&\\%_2'(",
					}, {
						entrance: 1901,
						id:       "1stDisplayId",
						nickname: " !\\%_1\"#",
						realname: "$&\\%_2'(",
					},
				],
				name: "Prog部",
			},
		}, {
			name:     "clubList",
			run:      () => client.clubList(),
			expected: [
				{
					chief: {
						id:       "2ndDisplayId",
						mail:     "",
						nickname: " !%_1\"#",
						realname: "$&\\%_2'(",
						tel:      "",
					},
					id:   "prog",
					name: "Prog部",
				},
			],
		}, {
			name:     "memberDetail",
			run:      () => client.memberDetail("1stDisplayId"),
			expected: {
				affiliation: "理学部第一部 数理情報科学科",
				clubs:       [
					{
						chief: false,
						id:    "prog",
						name:  "Prog部",
					},
				],
				entrance:  1901,
				gender:    "女",
				mail:      "1st@kagucho.net",
				nickname:  " !\\%_1\"#",
				ob:        false,
				positions: [
					{id: "president", name: "局長"},
				],
				realname: "$&\\%_2'(",
				tel:      "012-345-567",
			},
		}, {
			name:     "memberList",
			run:      () => client.memberList(),
			expected: {
				affiliation: "理学部第一部 数理情報科学科",
				clubs:       [
					{
						chief: false,
						id:    "prog",
						name:  "Prog部",
					},
				],
				entrance:  1901,
				gender:    "女",
				mail:      "1st@kagucho.net",
				nickname:  " !\\%_1\"#",
				ob:        false,
				positions: [{id: "president", name: "局長"}],
				realname:  "$&\\%_2'(",
				tel:       "012-345-567",
			},
		}, {
			name:     "officerDetail",
			run:      () => client.officerDetail("president"),
			expected: {
				member: {
					id:       "1stDisplayId",
					mail:     "1st@kagucho.net",
					nickname: " !\\%_1\"#",
					realname: "$&\\%_2'(",
					tel:      "012-345-567",
				},
				name:  "局長",
				scope: ["management", "privacy"],
			},
		}, {
			name:     "officerList",
			run:      () => client.officerList(),
			expected: [
				{
					id:     "president",
					member: {
						id:       "1stDisplayId",
						mail:     "1st@kagucho.net",
						nickname: " !\\%_1\"#",
						realname: "$&\\%_2'(",
						tel:      "012-345-567",
					},
					name: "局長",
				},
			],
		},
	];

	for (const test of staticTests) {
		(function() {
			const expected = JSON.stringify(this.expected);

			describe(this.name, () => {
				const staticPromise = signinPromise.then(this.run);

				it("should return a successful promise", function(done) {
					signinPromise.then(
						() => staticPromise.then(
							() => done(), done),
						this.skip.bind(this));
				});

				it("should return a promise resolved with " + expected, function(done) {
					staticPromise.then(data => {
						const result = JSON.stringify(data);
						done(result == expected ? null : "got " + result);
					}, this.skip.bind(this));
				});
			});
		}.bind(test))();
	}
});

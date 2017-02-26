/**
	@file api_test.js implements a testing code for api.js.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client/api_test */

import * as api from "./api";

describe("module:api", () => {
	const tokenPromise = api.getTokenWithPassword("1stDisplayID", "1stPassword");

	describe("getTokenWithPassword", () => {
		it("should return a success promise",
			done => void tokenPromise.then(() => done(), done));

		it("should notify progress",
			done => void tokenPromise.progress(
				progress => progress == 1 ? done() : null));

		it("should return access_token", function(done) {
			tokenPromise.then(
				data => done(data.access_token ? null :
					new Error("got " + data.access_token)),
				this.skip.bind(this));
		});

		it("should return refresh_token", function(done) {
			tokenPromise.then(
				data => done(data.refresh_token ? null :
					new Error("got " + data.refresh_token)),
				this.skip.bind(this));
		});
	});

	describe("getTokenWithRefreshToken", () => {
		const refreshPromise = tokenPromise.done(
			token => api.getTokenWithRefreshToken(token.refresh_token));

		it("should return a successful promise", function(done) {
			tokenPromise.then(
				() => refreshPromise.then(() => done(), done),
				this.skip.bind(this));
		});

		it("should notify progress", function(done) {
			tokenPromise.then(
				() => refreshPromise.progress(
					progress => progress == 1 ? done() : null),
				this.skip.bind(this));
		});

		it("should resolve with access_token", function(done) {
			refreshPromise.then(
				data => done(data.access_token ?
					null : new Error("invalid access_token; got " + data.access_token)),
				this.skip.bind(this));
		});
	});

	describe("memberUpdate", () => it("should reject with invalid entrance", function(done) {
		tokenPromise.then(token => {
			api.memberUpdate(token.access_token, {entrance: "0"}).then(
			done.bind(undefined, new Error("unexpected resolution")),
			xhr => {
				const expected = "value out of range";
				done(xhr.responseJSON.error_description != expected &&
					new Error("expected \""+xhr.responseJSON.error_description+
						"\", got \""+expected+"\""));
			});
		}, this.skip.bind(this));
	}));

	const staticConsumers = [
		{
			name:     "clubDetail",
			consume:  token => api.clubDetail(token.access_token, "prog"),
			expected: {
				chief: {
					id:       "2ndDisplayId", mail:     "",
					nickname: " !%_1\"#", realname: "$&\\%_2'(",
					tel:      "",
				},
				members: [
					{
						entrance: 1901, id:       "2ndDisplayId",
						nickname: " !%_1\"#", realname: "$&\\%_2'(",
					}, {
						entrance: 1901, id:       "1stDisplayId",
						nickname: " !\\%_1\"#", realname: "$&\\%_2'(",
					},
				],
				name: "Prog部",
			},
		}, {
			name:     "clubList",
			consume:  token => api.clubList(token.access_token),
			expected: [
				{
					chief: {
						id:       "2ndDisplayId", mail:     "",
						nickname: " !%_1\"#", realname: "$&\\%_2'(",
						tel:      "",
					},
					id:   "prog",
					name: "Prog部",
				},
			],
		}, {
			name:     "memberDetail",
			consume:  token => api.memberDetail(token.access_token, "1stDisplayId"),
			expected: {
				affiliation: "理学部第一部 数理情報科学科",
				clubs:       [{chief: false, id: "prog", name: "Prog部"}],
				entrance:    1901,
				gender:      "女",
				mail:        "1st@kagucho.net",
				nickname:    " !\\%_1\"#",
				ob:          false,
				positions:   [{id: "president", name: "局長"}],
				realname:    "$&\\%_2'(",
				tel:         "012-345-567",
			},
		}, {
			name:     "memberList",
			consume:  token => api.memberList(token.access_token),
			expected: {
				affiliation: "理学部第一部 数理情報科学科",
				clubs:       [{chief: false, id: "prog", name: "Prog部"}],
				entrance:    1901,
				gender:      "女",
				mail:        "1st@kagucho.net",
				nickname:    " !\\%_1\"#",
				ob:          false,
				positions:   [{id: "president", name: "局長"}],
				realname:    "$&\\%_2'(",
				tel:         "012-345-567",
			},
		}, {
			name:     "officerDetail",
			consume:  token => api.officerDetail(token.access_token, "president"),
			expected: {
				member: {
					id:       "1stDisplayId", mail:     "1st@kagucho.net",
					nickname: " !\\%_1\"#", realname: "$&\\%_2'(",
					tel:      "012-345-567",
				},
				name:  "局長",
				scope: ["management", "privacy"],
			},
		}, {
			name:     "officerList",
			consume:  token => api.officerList(token.access_token),
			expected: [
				{
					id:     "president",
					member: {
						id:       "1stDisplayId", mail:     "1st@kagucho.net",
						nickname: " !\\%_1\"#", realname: "$&\\%_2'(",
						tel:      "012-345-567",
					},
					name: "局長",
				},
			],
		},
	];

	for (const consumer of staticConsumers) {
		(function() {
			const expected = JSON.stringify(this.expected);

			describe(this.name, () => {
				const staticPromise = tokenPromise.then(this.consume);

				it("should return a successful promise", function(done) {
					tokenPromise.then(
						() => staticPromise.then(
							() => done(), done),
						this.skip.bind(this));
				});

				it("should notify progress", function(done) {
					tokenPromise.then(
						() => staticPromise.progress(
							progress => progress == 1 ? done() : null),
						this.skip.bind(this));
				});

				it("should return a promise resolved with " + expected, function(done) {
					staticPromise.then(data => {
						const result = JSON.stringify(data);
						done(result == expected ? null : "got " + result);
					}, this.skip.bind(this));
				});
			});
		}.bind(consumer))();
	}
});

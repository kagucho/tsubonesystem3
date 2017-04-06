/**
	@file client.js provides the integrated feature of the client.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client/client */

import * as api from "./api";
import newSession from "./session";
import Stream from "mithril/stream";

/**
	TODO
*/
function *objectEntries(object) {
	for (const key in object) {
		yield [key, object[key]];
	}
}

/**
	TODO
*/
function bindSyncMap(m, mKey, mContainer, n, nKey, nContainer) {
	m.bindings.set(mKey, {
		Container:         mContainer,
		referent:          n,
		ReferentContainer: nContainer,
		referentKey:       nKey,
	});

	n.bindings.set(nKey, {
		Container:         nContainer,
		referent:          m,
		ReferentContainer: mContainer,
		referentKey:       mKey,
	});
}

/**
	TODO
*/
function addReference(container, referentPatch, ownerID, value, {referent, ReferentContainer, referentKey}) {
	function handleReferent(referentID) {
		const entry = referent.map.get(referentID);

		if (ReferentContainer) {
			let referentContainer;

			switch (entry) {
			case null:
				return false;

			case undefined:
				referentContainer = new ReferentContainer;
				referent.map.set(referentID,
					new Map([{key: referentKey, value}]));
				break;

			default:
				referentContainer = entry.get(referentKey);
				if (!referentContainer) {
					referentContainer = new ReferentContainer;
					entry.set(referentKey, referentContainer);
				}
			}

			(referentContainer.add || referentContainer.set).call(referentContainer, ownerID);
		} else {
			switch (entry) {
			case null:
				return false;

			case undefined:
				referent.map.set(referentID,
					new Map([{key: referentKey, ownerID}]));
				break;

			default:
				entry.set(referentKey, ownerID);
			}
		}

		if (referentPatch) {
			let patchEntry = referentPatch.get(referentID);
			if (!patchEntry) {
				patchEntry = new Set;
				referentPatch.set(referentID, patchEntry);
			}

			patchEntry.add(referentKey);
		}

		return true;
	}

	if (container) {
		const iterator = value[Symbol.iterator] ?
			value[Symbol.iterator]() : objectEntries(value);

		if (container.add) {
			for (const entry of iterator) {
				if (handleReferent(entry)) {
					container.add(entry);
				}
			}
		} else if (container.set) {
			for (const [id, entry] of iterator) {
				if (handleReferent(id)) {
					container.set(id, entry);
				}
			}
		} else {
			throw new Error("unknown type of reference container");
		}

		return null;
	} else {
		return value;
	}
}

/**
	TODO
*/
function createReference(referentPatch, ownerID, value, binding) {
	let addedValue;
	let container;

	if (binding.Container) {
		container = new binding.Container;
	}

	return addReference(container, referentPatch, ownerID, value, binding) || container;
}

/**
	TODO
*/
function patchReference(reference, referentPatch, ownerID, value, binding) {
	function deleteReferentReference(id) {
		if (binding.referentContainer) {
			binding.referent.list.get(id).get(binding.referentKey).delete(ownerID);
		} else {
			binding.referent.map.get(reference).set(binding.referentKey, null);
		}

		let entry = referentPatch.get(id);
		if (!entry) {
			entry = new Set;
			referentPatch.set(id, entry);
		}

		entry.add(binding.referentKey);
	}

	let container;
	let newValue;

	if (binding.Container) {
		newValue = new reference.constructor(
			value[Symbol.iterator] ? value : objectEntries(value));

		if (newValue.get) {
			for (const [id, entry] of reference) {
				const newEntry = newValue.get(id);

				switch (newEntry) {
				case undefined:
					reference.delete(id);
					deleteReferentReference(id);
					break;

				case entry:
					newValue.delete(id);
					break;

				default:
					reference.set(id, newEntry);
					newValue.delete(id);
				}
			}
		} else {
			for (const id of reference) {
				if (newValue.delete(id)) {
					continue;
				}

				reference.delete(id);
				deleteReferentReference(id);
			}
		}

		container = reference;
	} else {
		deleteReferentReference(id);
		newValue = value;
	}

	return addReference(container, referentPatch, ownerID, newValue, binding);
}

/**
	TODO
*/
class SyncMap {
	/**
		TODO
	*/
	constructor(id) {
		this.bindings = new Map;
		this.id = id;
		this.map = new Map;
		this.listeners = $.Callbacks("unique");
	}

	/**
		TODO
	*/
	mapProperty(key, value) {
		const binding = this.bindings.get(key);

		return binding && binding.Container ?
			Object.freeze({
				entries: value.entries && function() {
					return value.entries(...arguments);
				},

				forEach: value.forEach && function() {
					return value.forEach(...arguments);
				},

				get: value.get && function() {
					return value.get(...arguments);
				},

				has: value.has && function() {
					return value.has(...arguments);
				},

				keys: value.keys && function() {
					return value.keys(...arguments);
				},

				values: value.values && function() {
					return value.values(...arguments);
				},

				[Symbol.iterator]: value[Symbol.iterator] && function() {
					return value[Symbol.iterator](...arguments);
				},
			}) :
			value;
	}

	/**
		TODO
	*/
	sync(id, remote) {
		let synced = this.map.get(id);

		switch (synced) {
		case null:
			break;

		case undefined:
			synced = new Map;
			this.map.set(id, synced);
			// fallthrough

		default:
			for (const key in remote) {
				const binding = this.bindings.get(key);
				if (binding) {
					synced.set(key, createReference(null, id, remote[key], binding));
				} else if (!synced.has(key)) {
					synced.set(key, remote[key]);
				}
			}
		}
	}

	/**
		TODO
	*/
	delete(fetch, id) {
		return fetch().done(() => {
			const entry = this.map.get(id);

			this.map.set(id, null);
			for (const [key, binding] of this.bindings) {
				const referentIDs = entry.get(key);

				if (!referentIDs) {
					continue;
				}

				for (const referentID of referentIDs) {
					const referentEntry = binding.referent.map.get(referentID);

					if (binding.referentContainer) {
						referentEntry.get(binding.referentKey).delete(id);
					} else {
						referentEntry.set(binding.referentKey, null);
					}
				}
			}

			this.listeners.fire(new Map([{key: id, value: null}]));
			for (const [key, binding] of this.bindings) {
				const referentIDs = entry.get(key);

				if (!referentIDs) {
					continue;
				}

				binding.referent.listeners.fire(new Map((function *() {
					for (const referentID of referentIDs) {
						yield {
							key:   referentID,
							binding: new Set([binding.referentKey]),
						};
					}
				})()));
			}
		});
	}

	/**
		TODO
	*/
	patch(fetch, id, properties) {
		return fetch().done(() => {
			let entry = this.map.get(id);

			if (!entry) {
				entry = new Map;
				this.map.set(id, entry);
			}

			const referentPatches = [];

			for (const key in properties) {
				const binding = this.bindings.get(key);
				if (binding) {
					const reference = entry.get(key);
					const patch = new Map;

					if (reference) {
						const newReference = patchReference(reference,
							patch, id, properties[key], binding);
						if (newReference) {
							entry.set(key, newReference);
						}
					} else {
						entry.set(key, createReference(patch, id, properties[key], binding));
					}

					referentPatches.push({
						referent: binding.referent,
						patch,
					});
				} else {
					entry.set(key, properties[key]);
				}
			}

			this.listeners.fire(
				new Map([
					{
						key:   id,
						value: new Set(Object.keys(properties)),
					},
				]));

			for (const {referent, patch} of referentPatches) {
				referent.listeners.fire(patch);
			}
		});
	}

	/**
		TODO
	*/
	detailMapper(fetch) {
		const streams = Object.create(null);

		return (id, callback) => {
			if (!streams[id]) {
				const map = () => {
					const entry = this.map.get(id);

					if (entry === null) {
						throw "not_found";
					}

					const mapped = {};

					for (const [key, value] of entry) {
						mapped[key] = this.mapProperty(key, value);
					}

					return Object.freeze(mapped);
				};

				streams[id] = Stream(
					fetch(id).then(detail => {
						this.sync(id, detail);

						this.listeners.add(
							patch => patch.has(id) &&
								streams[id]($.Deferred().resolve(map()).promise()));

						return map();
					}));
			}

			return streams[id].map(callback);
		};
	}

	/**
		TODO
	*/
	listMapper(fetch, keys) {
		const map = this;
		let stream;

		function *generator() {
			for (const [id, entry] of map.map) {
				if (!entry) {
					continue;
				}

				const mapped = {[map.id]: id};

				for (const key of keys) {
					mapped[key] = map.mapProperty(key, entry.get(key));
				}

				yield mapped;
			}
		}

		return callback => {
			if (!stream) {
				stream = Stream(fetch().then(entries => {
					for (const entry of entries) {
						const id = entry[this.id];
						delete entry[this.id];
						this.sync(id, entry);
					}

					this.listeners.add(patch => {
						for (const value of patch.values()) {
							if (!value || keys.some(key => value.has(key))) {
								stream($.Deferred().resolve(generator).promise());
								break;
							}
						}
					});

					return generator;
				}));
			}

			return stream.map(callback);
		};
	}
}

/**
	TODO
*/
class MemberMap extends SyncMap {
	/**
		TODO
	*/
	constructor() {
		super("id");
	}

	/**
		TODO
	*/
	detailMapper(fetch, clubs) {
		const streams = Object.create(null);

		return (id, requestID, callback) => {
			if (!streams[id]) {
				const map = () => {
					const entry = this.map.get(id);

					if (entry === null) {
						throw "not_found";
					}

					const mapped = {};

					for (const [key, value] of entry) {
						mapped[key] = key == "clubs" ?
							function *() {
								for (const clubID of value) {
									yield {
										id:    clubID,
										chief: clubs.map.get(clubID).chief == id,
									};
								}
							} : this.mapProperty(key, value);
					}

					return Object.freeze(mapped);
				};

				streams[id] = Stream(
					fetch(requestID).then(detail => {
						const rawClubs = detail.clubs;
						const fire = () => streams[id]($.Deferred().resolve(map()).promise());

						detail.clubs = {
							[Symbol.iterator]: function *() {
								for (const club of rawClubs) {
									yield club.id;
								}
							},
						},

						this.sync(id, detail);

						for (const club of rawClubs) {
							if (club.chief) {
								clubs.sync(club.id, {chief: id});
							}
						}

						this.listeners.add(
							patch => patch.has(id) && fire());

						clubs.listeners.add(patch => {
							for (const clubID of this.map.get(id).clubs) {
								if (patch.get(clubID).chief) {
									fire();
								}
							}
						});

						return map();
					}));
			}

			return streams[id].map(callback);
		};
	}
}

/**
	TODO
*/
class PartyMap extends SyncMap {
	/**
		TODO
	*/
	constructor() {
		super("name");
	}

	/**
		TODO
	*/
	userListMapper(getID, fetch) {
		const keys = ["creator", "start", "end", "place", "inviteds", "due"];
		const map = this;
		let stream;

		function *generator() {
			for (const [id, entry] of map.map) {
				if (!entry) {
					continue;
				}

				const mapped = {
					[map.id]: id,
					user:     entry.get("attendances").get(getID()),
				};

				for (const key of keys) {
					mapped[key] = map.mapProperty(key, entry.get(key));
				}

				yield mapped;
			}
		}

		return callback => {
			if (!stream) {
				stream = Stream(fetch().then(entries => {
					for (const entry of entries) {
						const id = entry[map.id];
						entry.attendances = {[getID()]: entry.user};
						delete entry[map.id];
						delete entry.user;
						map.sync(id, entry);
					}

					map.listeners.add(patch => {
						for (const value of patch.values()) {
							if (!value || keys.some(key => value.has(key)) || value.has("attendances")) {
								stream($.Deferred().resolve(generator).promise());
								break;
							}
						}
					});

					return generator;
				}));
			}

			return stream.map(callback);
		};
	}
}

/**
	TODO
*/
export function merge(...streams) {
	const results = new Map;
	let lengthComputable = false;
	let loaded = 0;
	let total = 0;
	let resolveds = 0;

	return Stream.combine(function() {
		const changeds = arguments[arguments.length - 1];
		const deferred = $.Deferred();

		for (const changed of changeds) {
			let result = results.get(changed);
			if (!result) {
				result = {loaded: 0, total: 0};
				results.set(changed, result);
			} else if (result.data) {
				result.data = null;
				resolveds--;
			}

			changed().then(data => {
				result.data = data;
				resolveds++;

				if (resolveds >= results.size) {
					deferred.resolve(...(function *() {
						for (const value of results.values()) {
							yield value.data;
						}
					})());
				}
			}, deferred.reject, function(event) {
				if (event.lengthComputable) {
					({lengthComputable} = event);
					loaded = event.loaded - result.loaded;
					total = event.total - result.loaded;

					result.loaded = event.loaded;
					result.total = event.total;
				}

				deferred.notify(
					{lengthComputable, loaded, total});
			});
		}

		return deferred;
	}, streams);
}

/**
	module:private/client/client is a class to provide the integrated
	feature of the client.
	@extends Object
*/
export default function() {
	const clubs = new SyncMap("id");
	const mails = new SyncMap("subject");
	const members = new MemberMap;
	const officers = new SyncMap("id");
	const parties = new PartyMap;
	const session = newSession();

	bindSyncMap(clubs, "members", Set, members, "clubs", Set);
	bindSyncMap(mails, "recipients", Set, members, "mails", Set);
	bindSyncMap(members, "positions", Set, officers, "member", null);
	bindSyncMap(parties, "attendances", Map, members, "parties", Set);

	const mapMemberPrimitive = members.detailMapper(
		id => session.applyToken(token => api.getMember(token, id)),
		clubs);

	return {
		/**
			clubDetail returns the details of the club identified with the given ID.
			TODO
			@param {!String} id - The ID.
			@returns {!module:private/promise} A promise resolved with the details.
		*/
		mapClub: clubs.detailMapper(
			id => session.applyToken(
				token => api.getClub(token, id))),

		/**
			clubs returns the clubs. TODO
			@returns {!module:private/promise} A promise resolved with the clubs.
		*/
		mapClubs: clubs.listMapper(
			session.applyToken.bind(session, api.getClubs),
			["name", "chief", "members"]),

		/**
			getFilling returns whether the user agent is prompting the user
			to fill his information.
			@returns {!Boolean} The boolean indicating whether the user
			agent is prompting the user to fill his information.
		*/
		getFilling() {
			return session.getFilling();
		},

		/**
			getID returns the ID of the member bound to the current session.
			@returns {!String} The ID.
		*/
		getID() {
			return session.getID();
		},

		/**
			getScope returns the scope of the current session.
			@returns {!String[]} The scope.
		*/
		getScope() {
			return session.getScope();
		},

		/**
			mailCreate creates a email. TODO
			@returns {!module:private/promise} A promise desribing the
			progress and the result.
		*/
		createMail(id, properties) {
			return mails.patch(
				session.applyToken.bind(session, token => api.putMail(token, id, properties)),
				id, properties);
		},

		/**
			TODO
		*/
		mapMail: mails.detailMapper(
			subject => session.applyToken(
				token => api.getMail(token, subject))),

		/**
			TODO
		*/
		mapMails: mails.listMapper(
			session.applyToken.bind(session, api.getMails),
			["date", "from", "to"]),

		/**
			memberCreate creates a member. TODO
			@param {!*} properties - The properties of the new member.
			@returns {!module:private/promise} A promis describing the
			progress and the result.
		*/
		createMember(id, properties) {
			return members.patch(
				session.applyToken.bind(session,
					token => api.putMember(token, id, properties)),
				id, properties);
		},

		/**
			memberDetail returns the details of the member identified with
			the given ID. TODO
			@param {!String} id - The ID.
			@returns {!module:private/promise} A promise resolved with the
			details.
		*/
		mapMember(id, callback) {
			return mapMemberPrimitive(id, id, callback);
		},

		/**
			memberDelete deletes a member identified with the given ID.
			TODO
			@param {!String} id - The ID.
			@returns {!module:private/promise} A promise describing the
			result.
			FIXME: reflect this to clubs, and so on.
		*/
		deleteMember(id) {
			return members.delete(
				session.applyToken.bind(session,
					token => api.deleteMember(token, id)),
				id);
		},

		/**
			memberList returns the members. TODO
			@returns {!module:private/promise} A promise resolved with the members.
		*/
		mapMembers: members.listMapper(
			session.applyToken.bind(session, api.getMembers),
			["affiliation", "entrance", "nickname", "ob", "realname"]),

		/**
			TODO
		*/
		mapMemberMails: members.listMapper(
			() => session.applyToken(api.getMembersMails).then(function *(map) {
				console.log(map);
				for (const id in map) {
					yield {id, mail: map[id]};
				}
			}), ["mail"]),

		/**
			officerDetail returns the details of the officer identified with
			the given ID. TODO
			@param {!String} id - The ID.
			@returns {!module:private/promise} A promise resolved with the
			details.
		*/
		mapOfficer: officers.detailMapper(
			id => session.applyToken(
				token => api.getOfficer(token, id))),

		/**
			officerList returns the officers. TODO
			@returns {!module:private/promise} A promise resolved with the officers.
		*/
		mapOfficers: officers.listMapper(
			session.applyToken.bind(session, api.getOfficers),
			["name", "member"]),

		/**
			partyCreate creates a party. TODO
			@param {!*} properties - Properties of the party.
			@returns {!module:private/promise} A promise describing the
			result.
		*/
		createParty(name, properties) {
			return parties.patch(
				session.applyToken.bind(session,
					token => api.putParty(token, name, properties)),
				name, $.extend({attendances: []}, properties));
		},

		/**
			partyList lists the parties. TODO
			@returns {!module:private/promise} A promise resolved with the
			parties.
		*/
		mapParties: parties.userListMapper(
			session.getID,
			session.applyToken.bind(session, api.getParties)),

		/**
			TODO
		*/
		respondParty(party, attending) {
			return parties.patch(
				session.applyToken.bind(session,
					token => api.patchParty(token, party, {attending: attending ? "1" : "0"})),
				party, {attendances: {[session.getID()]: attending ? "accepted" : "declined"}});
		},

		/**
			recoverSession recovers the session from sessionStorage. TODO
			@returns {!module:private/promise} A promise resolved when
			recovered.
		*/
		recoverSession() {
			return session.recover();
		},

		/**
			setFillingToken sets the token to fill the information of the
			member. TODO
			@param {!String} id - The ID of the user.
			@param {!String} token - The access token.
			@returns {Undefined}
		*/
		setFillingToken(id, token) {
			session.setFillingToken(id, token);
		},

		/**
			signin signs in. TODO
			@param {!String} id - The ID.
			@param {!String} password - The password.
			@returns {!module:private/promise} A promise describing the
			progress and the result.
		*/
		signin(id, password) {
			return session.signin(id, password);
		},

		/**
			userConfirm confirms the email address of the user by submitting
			the token sent to the address. TODO
			@param {!String} mailToken - The token sent to the address.
			@returns {!module:private/promise} A promise describing the
			progress and the result.
		*/
		confirm(mailToken) {
			return members.patch(
				session.applyToken.bind(session,
					token => api.patchMember(token, null, {confirm: mailToken})),
				session.getID(), {confirmed: true});
		},

		/**
			userDeclareOB declares the user is OB. TODO
			@returns {!module:private/promise} A promise describing the
			progress and the result.
		*/
		declareOB() {
			return members.patch(
				session.applyToken.bind(session,
					token => api.patchMember(token, null, {ob: true})),
				session.getID(), {ob: true});
		},

		/**
			userDetail returns the details of the user. TODO
			@returns {!module:private/promise} A promise resolved with the
			details.
		*/
		mapUser(callback) {
			return mapMemberPrimitive(session.getID(), null, callback);
		},

		/**
			userUpdate updates properties of the user. TODO
			@param {!*} properties - The properties to update.
			@returns {!module:private/promise} A promise describing the
			progress and the result.
		*/
		patchUser(properties) {
			const remote = $.extend({}, properties);
			const local = {};

			if (remote.clubs) {
				remote.clubs = remote.clubs.join(" ");
			}

			if (properties.mail) {
				local.confirmed = false;
			}

			return members.patch(
				session.applyToken.bind(session,
					token => api.patchMember(token, null, remote).then(({scope, refresh_token, access_token}) => {
						if (scope && refresh_token) {
							session.updateToken(scope,
								refresh_token,
								access_token);
						}
					})),
				session.getID(), $.extend(local, properties));
		},
	};
}

/**
 * @file index.js implements the entry point for the private page.
 * @author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
 * @copyright Kagucho 2016
 * @license AGPL-3.0
 */

/** @module private */

import "./tags/app/club.tag";
import "./tags/app/clubs.tag";
import "./tags/app/index.tag";
import "./tags/app/member.tag";
import "./tags/app/members.tag";
import "./tags/app/officer.tag";
import "./tags/app/officers.tag";
import "./tags/table/member.tag";
import "./tags/table/officer.tag";
import "./tags/container.tag";
import "./tags/notfound.tag";
import "./tags/recover-session.tag";
import "./tags/signin.tag";
import "./tags/top-progress.tag";
import Client from "./client";

/**
 * client is a module:client.
 * @type !module:client
 */
const client = new Client("/");

const recoverDeferred = $.Deferred();
riot.mount("#container", "recover-session",
           {client, deferred: recoverDeferred});
recoverDeferred.catch(() => {
  const signinDeferred = $.Deferred();
  riot.mount("#container", "signin", {client, deferred: signinDeferred});

  return signinDeferred;
}).then(() => route.start(true));

route.base("/private");

route("#!clubs..", () => riot.mount("#container", "app-clubs", {client}));
route("#!club..", () => riot.mount("#container", "app-club", {client}));

route("#!members..",
           () => riot.mount("#container", "app-members", {client}));

route("#!member..",
           () => riot.mount("#container", "app-member", {client}));

route("#!officers..",
           () => riot.mount("#container", "app-officers", {client}));

route("#!officer..",
           () => riot.mount("#container", "app-officer", {client}));

route("", () => riot.mount("#container", "app-index"));
route("..", () => riot.mount("#container", "notfound"));

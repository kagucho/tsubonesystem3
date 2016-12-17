<!--
  Copyright (C) 2016  Kagucho <kagucho.net@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->
<app-members>
  <top-progress ratio={privateProgress} />
  <container>
    <div class="container">
      <div>
        <h1>Members</h1>
      </div>
      <div>
        <div class="col-sm-12 form-horizontal"
             style="background-color: #f0f8ff; border-radius: 1em; padding: 1em;">
          <div class="form-group">
            <label class="col-sm-3 control-label" for="yieldedNickname">
              ニックネーム
            </label>
            <div class="col-sm-9">
              <input class="form-control" maxlength="64" ref="yieldedNickname"
                     oninput={parent.privateOnInput} placeholder="Nickname">
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-3 control-label" for="yieldedRealname">
                名前
              </label>
            <div class="col-sm-9">
              <input class="form-control" maxlength="64" ref="yieldedRealname"
                     oninput={parent.privateOnInput} placeholder="Name">
            </div>
          </div>
          <div class="form-group">
            <label class="col-sm-3 control-label" for="yieldedEntrance">
              入学年度
            </label>
            <div class="col-sm-9">
              <input class="form-control" max="2155" min="1901"
                     ref="yieldedEntrance" oninput={parent.privateOnInput}
                     placeholder="Entrance" type="number">
            </div>
          </div>
          <div class="form-group">
            <label class="col-sm-3 control-label" for="yieldedOB">
              OB
            </label>
            <div class="col-sm-9">
              <input ref="yieldedOB" oninput={parent.privateOnInput}
                     type="checkbox">
            </div>
          </div>
          <div class="form-group">
            <label class="col-sm-3 control-label" for="yieldedActive">
              現役
            </label>
            <div class="col-sm-9">
              <input ref="yieldedActive" oninput={parent.privateOnInput}
                     type="checkbox">
            </div>
          </div>
        </div>
        <div style="clear: both; padding-top: 1em;">
          <p class="lead" if={parent.privateMembersMatchedNumber != null}
             style="color: gray;">
            {parent.privateMembersMatchedNumber} 件
          </p>
          <div if={parent.privateError} class="alert alert-danger"
               role="alert">
            <span class="glyphicon glyphicon-exclamation-sign"
                  aria-hidden="true"></span>
              {parent.privateError}
          </div>
        </div>
        <table class="table table-responsive">
          <thead>
            <tr style="background-color: #d9edf7;">
              <th>ニックネーム</th>
              <th>名前</th>
              <th>入学年度</th>
            </tr>
          </thead>
          <tbody>
            <tr each={parent.privateMembers} if={matched}>
              <td><a href="#!member?id={id}">{nickname}</a></td>
              <td>{realname}</td>
              <td>{entrance}</td>
            </tr>
         </tbody>
        </table>
      </div>
    </div>
  </container>
  <script>
    const fail = error => {
      this.privateError = error || "どうしようもないエラーが発生しました。";
      delete this.privateProgress;
      this.update();
    };

    try {
      const memberList = opts.client.memberList().then(data => {
        this.privateMembers = data;
      }, xhr => {
        let message;

        switch (xhr.status) {
        case 0:
          message = "TsuboneSystemへの経路上に問題が発生しました。ネットワーク接続などを確認してください。";
          break;

        case 401:
          message = "あんた誰?って言われちゃいました。もう一度サインインしてください。";
          break;

        default:
          message = "サーバー側のエラーです。がびーん。";
        }

        fail(message);
      }, progress => {
        this.privateProgress = progress;
        this.update();
      }).catch(() => fail());

      this.privateProgress = 0;

      this.on("mount", () => {
        const containerRefs = this.tags.container.refs;

        const queryMember = () => {
          this.privateMembersMatchedNumber = 0;

          for (const member of this.privateMembers) {
            member.matched =
              member.nickname.includes(containerRefs.yieldedNickname.value) &&
              member.realname.includes(containerRefs.yieldedRealname.value) &&
              (!containerRefs.yieldedEntrance.value ||
               member.entrance == containerRefs.yieldedEntrance.value) &&
              ((containerRefs.yieldedOB.checked && member.ob) ||
               (containerRefs.yieldedActive.checked && !member.ob));

            if (member.matched)
              this.privateMembersMatchedNumber++;
          }

          this.update();
        };

        this.privateOnInput = () => {
          try {
            const query = {};

            if (containerRefs.yieldedEntrance.value)
              query.entrance = containerRefs.yieldedEntrance.value;

            if (containerRefs.yieldedNickname.value)
              query.nickname = containerRefs.yieldedNickname.value;

            if (containerRefs.yieldedRealname.value)
              query.realname = containerRefs.yieldedRealname.value;

            if (containerRefs.yieldedOB.checked &&
                !containerRefs.yieldedActive.checked)
              query.ob = "1";
            else if (!containerRefs.yieldedOB.checked &&
                     containerRefs.yieldedActive.checked)
              query.ob = "0";

            let uri = "#!members";

            if (!$.isEmptyObject(query))
              uri = [uri, $.param(query)].join("?");

            history.pushState({}, document.title, uri);

            if (this.privateMembers)
              queryMember();
          } catch (exception) {
            fail();
          }
        };

        const query = route.query();

        if (query.nickname)
          containerRefs.yieldedNickname.value = query.nickname;

        if (query.realname)
          containerRefs.yieldedRealname.value = query.realname;

        if (query.entrance)
          containerRefs.yieldedEntrance.value = query.entrance;

        switch (query.OB) {
        case "0":
          containerRefs.yieldedActive.checked = true;
          break;

        case "1":
          containerRefs.yieldedOB.checked = true;
          break;

        default:
          containerRefs.yieldedOB.checked = true;
          containerRefs.yieldedActive.checked = true;
        }

        memberList.then(() => queryMember());
      });
    } catch (exception) {
      console.log(exception);
      fail();
    }
  </script>
</app-members>

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
<app-member>
  <top-progress ratio={privateProgress} />
  <container>
    <div class="container">
      <div if={parent.privateError} class="alert alert-danger" role="alert">
        <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true">
        </span>
        {parent.privateError}
      </div>
      <div if={parent.privateMember}>
        <h1 if={parent.privateMember.nickname != null}
            style="font-size: x-large;">
          {parent.privateMember.nickname}ちゃんの詳細情報
        </h1>
        <table class="table">
          <tr>
            <th>ID</th>
            <td if={parent.privateMember.id != null}>
              {parent.privateMember.id}
            </td>
            <td if={parent.privateMember.id == null}>?</td>
          </tr>
          <tr>
            <th>ニックネーム</th>
            <td if={parent.privateMember.nickname != null}>
              {parent.privateMember.nickname}
            </td>
            <td if={parent.privateMember.nickname == null}>?</td>
          </tr>
          <tr>
            <th>名前</th>
            <td if={parent.privateMember.realname != null}>
              {parent.privateMember.realname}
            </td>
            <td if={parent.privateMember.realname == null}>?</td>
          </tr>
          <tr>
            <th>性別</th>
            <td if={parent.privateMember.gender != null}>
              {parent.privateMember.gender}
            </td>
            <td if={parent.privateMember.gender == null}>?</td>
          </tr>
          <tr>
            <th>メールアドレス</th>
            <td if={parent.privateMember.mail != null}>
              <a href="mailto:{encodeURIComponent(parent.privateMember.mail)}">
                {parent.privateMember.mail}
              </a>
            </td>
            <td if={parent.privateMember.mail == null}>?</td>
          </tr>
          <tr>
            <th>電話番号</th>
            <td if={parent.privateMember.tel != null}>
              <a href="tel:{parent.privateMember.tel.replace(/^0(?!-)/, "+81-")}">
                {parent.privateMember.tel}
              </a>
            </td>
            <td if={parent.privateMember.tel == null}>?</td>
          </tr>
          <tr>
            <th>役職</th>
            <td if={parent.privateMember.positions != null}>
              <div each={parent.privateMember.positions}>
                <a href="#!officer?id={id}">{name}</a>
              </div>
            </td>
            <td if={parent.privateMember.positions == null}>?</td>
          </tr>
          <tr>
            <th>所属部</th>
            <td if={parent.privateMember.clubs != null}>
              <div each={parent.privateMember.clubs}>
                <a href="#!club?id={id}">{name}</a>
                <span if={chief}>(部長)</span>
              </div>
            </td>
            <td if={parent.privateMember.clubs == null}>?</td>
          </tr>
          <tr>
            <th>学科</th>
            <td if={parent.privateMember.affiliation != null}>
              {parent.privateMember.affiliation}
            </td>
            <td if={parent.privateMember.affiliation == null}>?</td>
          </tr>
          <tr>
            <th>入学年度</th>
            <td if={parent.privateMember.entrance != null}>
              {parent.privateMember.entrance}
            </td>
            <td if={parent.privateMember.entrance == null}>?</td>
          </tr>
          <tr>
            <th>OB宣言</th>
            <td if={parent.privateMember.ob == true}>OB宣言済み</td>
            <td if={parent.privateMember.ob == false}>(現役部員)</td>
            <td if={parent.privateMember.ob == null}>?</td>
          </tr>
          <style scoped>
            th {
              font-size: 1.25em;
              font-weight: 500;
            }
  
            td, th {
              padding: 0.75em !important;
              vertical-align: middle !important;
            }
          </style>
        </table>
      </div>
    </div>
  </container>
  <script>
    const fail = error => {
      this.privateError = error ? error : "どうしようもないエラーが発生しました。";
      delete this.privateProgress;
      this.update();
    };

    try {
      const failXHR = xhr => {
        let message;

        switch (xhr.status) {
        case 400:
          if (xhr.responseJSON.error == "invalid_id") {
            message = "IDが違うってよ。";
            break;
          }
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
      };

      const query = route.query();
      this.privateMember = {id: query.id};

      const memberPromise = opts.client.memberDetail(query.id).then(data => {
        Object.assign(this.privateMember, data);
        this.update();
      }, failXHR, progress => {
        this.privateProgress = progress;
        this.update();
      }).catch(() => fail());

      this.privateProgress = 0;
    } catch (exception) {
      fail();
    }
  </script>
</app-member>

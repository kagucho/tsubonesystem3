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
<app-officer>
  <top-progress ratio={privateProgress} />
  <container>
    <div class="container">
      <div if={parent.privateError} class="alert alert-danger" role="alert">
        <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true">
        </span>
        {parent.privateError}
      </div>
      <div if={parent.privateOfficer}>
        <div>
          <h1 style="font-size: x-large;">
            {parent.privateOfficer.name}閣下の詳細情報
          </h1>
        </div>
        <div>
          <h2 style="font-size: large;">権限</h2>
          <ul>
            <li each={parent.privateOfficer.scope}>{description}</li>
          </ul>
        </div>
        <div>
          <table-officer items={[parent.privateOfficer.member]} />
        </div>
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
      opts.client.officerDetail(route.query().id).then(data => {
        const scopeDescription = {
          management: "メンバー情報を更新できる",
          privacy: "メンバーの電話番号を閲覧できる",
        };

        data.scope.map((id, index, array) => {
          array[index] = {description: scopeDescription[id]};
        });

        this.privateOfficer = data;
        this.update();
      }, xhr => {
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
      }, progress => {
        this.privateProgress = progress;
        this.update();
      }).catch(() => fail());

      this.privateProgress = 0;
    } catch (exception) {
      fail();
    }
  </script>
</app-officer>

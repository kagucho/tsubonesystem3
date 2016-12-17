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
<app-club>
  <top-progress ratio={privateProgress} />
  <container>
    <div class="container">
      <div if={parent.privateError} class="alert alert-danger" role="alert">
        <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true">
        </span>
        {parent.privateError}
      </div>
      <div if={parent.privateClub}>
        <div>
          <h1>{parent.privateClub.name}の詳細情報</h1>
        </div>
        <div>
          <h2>部長</h2>
          <table-officer items={[parent.privateClub.chief]} />
        </div>
        <div>
          <h2>{parent.privateClub.name}のいかれた仲間たち</h2>
          <p if={parent.privateClub.members.length} style="color: gray;">
            {parent.privateClub.members.length} 件
          </p>
          <table-member items={parent.privateClub.members} />
        </div>
        <style scoped>
          h1 {
            font-size: x-large;
          }

          h2 {
            font-size: large;
          }
        </style>
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
      opts.client.clubDetail(route.query().id).then(data => {
        this.privateClub = data;
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
</app-club>

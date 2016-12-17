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
<app-clubs>
  <top-progress ratio={internal.progress} />
  <container>
    <div class="container">
      <div>
        <h1>Clubs</h1>
      </div>
      <div if={parent.internal.error} class="alert alert-danger" role="alert">
        <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true">
        </span>
        {parent.internal.error}
      </div>
      <div>
        <div each={parent.internal.clubs}>
          <h2>{name}</h2>
          <div>
            <h3>部長</h3>
            <table-officer items={[chief]} />
          </div>
          <div>
            <h3>
              <a href="#!club?id={id}">
                {name}のいかれた仲間たちを見る
              </a>
            </h3>
          </div>
        </div>
        <style scoped>
          h2 {
            font-size: x-large;
          }

          h3 {
            font-size: large;
          }
        </style>
      </div>
    </div>
  </container>
  <script>
    this.internal = {progress: 0};

    const fail = error => {
      this.internal.error = error ? error : "どうしようもないエラーが発生しました。";
      delete this.internal.progress;
      this.update();
    };

    try {
      opts.client.clubList().then(data => {
        this.internal.clubs = data;
        this.update();
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
        this.internal.progress = progress;
        this.update();
      }).catch(() => fail());
    } catch (exception) {
      fail();
    }
  </script>
</app-clubs>

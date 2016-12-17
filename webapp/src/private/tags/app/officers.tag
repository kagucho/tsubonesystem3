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
<app-officers>
  <top-progress ratio={privateProgress.sum} />
  <container>
    <div class="container">
      <h1>Officers</h1>
      <div if={parent.privateError} class="alert alert-danger" role="alert">
        <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true">
        </span>
        {parent.privateError}
      </div>
      <div>
        <div each={parent.privateOfficers}>
          <h2 class="text-center" style="font-size: x-large;">
            <a href="#!officer?id={id}">{name}</a>
          </h2>
          <table-officer items={[member]} />
        </div>
        <div>
          <h2 class="text-center" style="font-size: x-large;">各部長</h2>
          <div each={parent.privateClubs}>
            <h3 class="lead" style="margin-bottom: 0;">
              <a href="#!club?id={id}">
                {name}
              </a>
            </h3>
            <table-officer items={[chief]} />
          </div>
        </div>
      </div>
    </div>
  </container>
  <script>
    const fail = (message) => {
      this.privateError = message ?
                          message : "どうしようもないエラーが発生しました。";
      delete this.privateProgress.sum;
      this.update();
    };

    try {
      const officers = opts.client.officerList().then(data => {
        this.privateOfficers = data;
        this.update();
      }, xhrFail, progress => {
        this.privateProgress.officers = progress;
        this.privateProgress.sum = (progress + this.privateProgress.clubs) / 2;
        this.update();
      }).catch(() => fail());

      opts.client.clubList().then(data => {
        this.privateClubs = data;
        this.update();
      }, xhrFail, progress => {
        this.privateProgress.clubs = progress;
        this.privateProgress.sum =
          (progress + this.privateProgress.officers) / 2;
        this.update();
      }).catch(() => fail());

      function xhrFail(xhr) {
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
      }

      this.privateProgress = {officers: 0, clubs: 0, sum: 0};
    } catch (exception) {
      fail();
    }
  </script>
</app-officers>

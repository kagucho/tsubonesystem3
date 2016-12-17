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
<signin>
  <top-progress ratio={privateProgress} />
  <div class="jumbotron"
       style="margin-top: 2em; margin-bottom: 2em; padding: 0;">
    <div class="container text-center">
      <h1>TsuboneSystem</h1>
      <p>TsuboneSystemは出欠席管理、メンバー管理、簡易メーリングリスト、非常時連絡先参照を目的に作られたシステムです。</p>
      <p class="hidden-xs"></p>
    </div>
  </div>
  <div class="container">
    <p style="color: red;">{privateError}</p>
    <div class="text-center" style="max-width: 24em; margin: 0 auto;">
      <input class="form-control" maxlength="64" placeholder="ID"
             ref="id" type="text" />
      <input class="form-control" maxlength="64" placeholder="Password"
             ref="password" type="password" />
      <input class="btn btn-lg btn-primary btn-block"
             disabled={privateProgress != null} onclick={privateSubmit}
             type="button" value="Sign in" />
      <style scoped>
        input {
          font-size: 1.5em;
          height: auto;
          margin-top: 1em;
        }
      </style>
    </div>
  </div>
  <script>
    const fail = error => {
      this.privateError = error ? error : "どうしようもないエラーが発生しました。";
      this.update();
    }

    try {
      this.privateSubmit = () => {
        const failXHR = error => {
          delete this.privateProgress;
          fail(error);
        }

        try {
          opts.client.signin(this.refs.id.value,
                             this.refs.password.value).then(
          () => opts.deferred.resolve(),
          xhr => {
            switch (xhr.status) {
            case 401:
              failXHR("残念！！IDもしくはパスワードが違います。");
              break;

            case 429:
              failXHR("残念！！やり直しすぎです。ちょっと待ってから再度入力してください。");
              break;

            default:
              failXHR("TsuboneSystemへの経路上に問題が発生しました。ネットワーク接続などを確認してください。");
            }
          }, progress => {
            this.privateProgress = progress;
            this.update();
          }).catch(() => failXHR());

          this.privateProgress = 0;
          this.update();
        } catch (exception) {
          failXHR();
        }
      };
    } catch (exception) {
      fail();
    }
  </script>
</signin>

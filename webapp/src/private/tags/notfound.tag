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
<notfound>
  <container>
    <img alt="ロゴ" class="hidden-xs pull-right" src={parent.privateLogo} />
    <h1>ページが見つかりません</h1>
    <h2>どうしよう?</h2>
    <h3>このサイトにあるリンクを踏んでここに来た。</h3>
    <p>
      <a href="mailto:kagucho.net@gmail.com">kagucho.net@gmail.com</a>
      へ報告してください。お願いします。あなたの好意が世界を救います。
    </p>
    <h3>自分でURLを入力した。</h3>
    <p>
      打ち間違えてないか確認してください。打ち間違えてない? それは困った。お役に立てそうにない。
    </p>
    <h3>トライフォースを探しにここに来た。</h3>
    <p>ここにもないですから！本当に。</p>
    <footer class="clearfix"
            style="background-image: url({parent.privateFooter}); background-position: right; background-repeat: no-repeat; background-size: contain; border-color: black; border-top-style: solid; border-width: 0.125em; padding-top: 1em;">
      <p>
        Copyright (C) 2016 神楽坂一丁目通信局. Licensed under
        <a href="https://www.gnu.org/licenses/agpl-3.0.en.html">AGPL-3.0</a>.
      </p>
    </footer>
  </container>
  <script>
    import footer from "../../images/footer.png";
    import logo from "../../images/logo250c_black.png";

    this.privateFooter = footer;
    this.privateLogo = logo;
  </script>
</notfound>

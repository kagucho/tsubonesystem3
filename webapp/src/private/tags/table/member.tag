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
<table-member>
  <table class="table table-responsive">
    <thead>
      <tr style="background-color: #d9edf7;">
        <th>ニックネーム</th>
        <th>名前</th>
        <th>入学年度</th>
      </tr>
    </thead>
    <tbody>
      <tr each={ opts.items }>
        <td><a href="#!member?id={ id }">{ nickname }</a></td>
        <td>{ realname }</td>
        <td>{ entrance }</td>
      </tr>
    </tbody>
  </table>
</table-member>

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
<container>
  <header class="navbar-fixed-top">
    <div class="navbar navbar-default"
         style="background-color: white; border-style: hidden;">
      <div class="container" >
        <div class="navbar-header">
          <button type="button" class="navbar-toggle"
                  data-toggle="collapse" data-target=".navbar-top">
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="navbar-brand" href="">TsuboneSystem</a>
        </div>
        <div class="collapse navbar-collapse navbar-top">
          <ul class="nav navbar-nav">
            <li><a href="#!mail">Mail</a></li>
          </ul>
          <ul class="nav navbar-nav">
            <li><a href="#!party">Party</a></li>
          </ul>
          <ul class="nav navbar-nav">
            <li><a href="#!members">Members</a></li>
          </ul>
          <ul class="nav navbar-nav">
            <li><a href="#!clubs">Clubs</a></li>
          </ul>
          <ul class="nav navbar-nav">
            <li><a href="#!officers">Officers</a></li>
          </ul>
          <ul class="nav navbar-nav navbar-right">
            <li><a href="#!signout">Sign out</a></li>
          </ul>
          <ul class="nav navbar-nav navbar-right">
            <li><a href="#!settings">Settings</a></li>
          </ul>
        </div>
      </div>
    </div>
  </header>
  <div id="container-app" style="padding-top: 70px;"><yield /></div>
</container>

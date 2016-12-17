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
<top-progress>
  <!-- The z-index of navbar is 1030, so set the one of top-progress bigger. -->
  <div class={animation: opts.ratio == 1} if={opts.ratio != null}
       style="background-color: lightgray; position: fixed; top: 0px; height: 0.25em; width: 100%; z-index: 1031;">
    <div style="background-color: dodgerblue; height: 100%; width: { opts.ratio * 100 }%;">
    </div>
  </div>
  <style scoped>
    .animation {
      animation: 0.125s ease 1s 1 normal forwards running animation;
    }

    @keyframes animation {
      from {
        height: 0.25em;
      }

      to {
        height: 0em;
      }
    }
  </style>
</top-progress>

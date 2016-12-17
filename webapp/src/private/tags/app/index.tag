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
<app-index>
  <container>
    <div ref="yieldedSlide">
      <div class="crossfade-element" each={slide in parent.privateSlides}>
        <div class="slide"
             style="background-image: url({slide}); z-index: inherit;">
        </div>
      </div>
    </div>
    <style scoped>
      .slide {
        background-size: cover;
        position: fixed;
        height: 100%;
        width: 100%;
      }
    </style>
  </container>
  <script>
    import Crossfade from "../../../crossfade";
    import sakuraDrawing from "./slides/sakura_drawing.jpg";
    import sakuraPhoto from "./slides/sakura_photo.jpg";
    import skyDrawing from "./slides/sky_drawing.jpg";

    this.privateSlides = [sakuraDrawing, sakuraPhoto, skyDrawing];

    this.on("mount", () => {
      this.privateCrossfade =
        new Crossfade(this.tags.container.refs.yieldedSlide);

      this.privateCrossfade.start();
    });

    this.on("before-unmount", () => this.privateCrossfade.stop());
  </script>
</app-index>

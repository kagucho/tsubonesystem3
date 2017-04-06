/**
	@file root.js implements root component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/root */

/**
	module:private/component/app/root is a component to show the top page
	and the entrance of the Lost Woods. (i.e. the code is really mess!)
	@name module:private/component/app/root
	@type !external:Mithril~Component
*/

import * as container from "../container";
import slides from "../../../slides";

/**
	sqrt3 is square root of 3. Someone in this club also named "root 3".
	@private
	@type !Number
*/
const sqrt3 = Math.sqrt(3);

/**
	fixeds is a set of the addresses of the fixed cells. A fixed cell is a
	cell which cannot be a destination or a transit of an unfixed cell's
	move.
	@private
	@type !module:private/component/app/root~Address[]
*/
const fixeds = [{col: 7, row: 2}, {col: 4, row: 5}, {col: 10, row: 4}];

/**
	initUnfixeds is a set of the initial addresses of the unfixed cells. A
	unfixed cell is a cell which can be moved by the user.
	@private
	@type !module:private/component/app/root~Address[]
*/
const initUnfixeds = [{col: 7, row: 0}, {col: 3, row: 5}, {col: 3, row: 6}];

/**
	triforce is a set of the addresses of the Triforce's cells. The Triforce
	is a sacred relic which lies on the greatest kingdom and the greatest
	"system".
	@private
	@type !module:private/component/app/root~Address[]
*/
const triforce = [
	{col: 7, row: 3}, // power
	{col: 6, row: 4}, // wisdom
	{col: 8, row: 4}, // courage
];

// FOR DEBUGGING (and cheating :P)
// const initUnfixeds = triforce;

/**
	lineScale is a coefficient for calculating the scale.
	@private
	@type !Number
*/
const lineScale = 10;

/**
	lineElapsedToScale is a coefficient for calculating the length of the
	line at a moment.
	@private
	@type !Number
*/
const lineElapsedToScale = 512;

/**
	cellIsHyrule returns whether the cell identified by the given address
	is of Hyrule or not. If not, it is of Lorule.
	@private
	@param {!module:private/component/app/root~Address} cell - The address
	of the cell.
	@returns {!Boolean} A boolean indicating whether the given cell is of
	Hyrule or not.
*/
function cellIsHyrule(cell) {
	return !((cell.col + cell.row) % 2);
}

/**
	cellToCoordinate returns the absolute coordinate of the cell identified
	by the given address. This function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@param {!module:private/component/app/root~Address} cell - The address
	of the cell, relative to the left top cell.
	@returns {!module:private/component/app/root~Coordinate} The coordinate
	of the cell. If it is of Hyrule, the coordinate points to the top. If it
	is of Lorule, the coordinate points to the bottom middle.
*/
function cellToCoordinate({col, row}) {
	return {
		x: this.cellOrigin.x + col * this.scaleWidth,
		y: this.cellOrigin.y + row * this.scaleHeight,
	};
}

/**
	coordinateToCell returns the address of the cell including the point
	identified by the given absolute coordinate. This function is usable
	after init.
	@private
	@this module:private/component/app/root~Context
	@param {!module:private/component/app/root~Coordinate} abs - The
	absolute coordinate identifying a point in the cell.
	@returns {!module:private/component/app/root~Address} The address of
	the cell, relative to the left top cell.
*/
function coordinateToCell({x: absX, y: absY}) {
	const x = absX - this.cellOrigin.x;
	const y = absY - this.cellOrigin.y;
	const xFloor = Math.floor(x / this.scaleWidth);
	const yFloor = Math.floor(y / this.scaleHeight);
	const yModWidth = (y % this.scaleHeight) / sqrt3;

	return {
		col: (absX - this.tileOrigin.x) % this.scaleWidth < ((xFloor + yFloor) % 2 ? this.scaleWidth - yModWidth : yModWidth) ?
			xFloor : xFloor + 1,

		row: yFloor || 0,
	};
}

/**
	fillCell fills the cell identified by the given address in the given
	context. This function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@param {!external:DOM~CanvasRenderingContext2D} context - The context
	where the cell will be filled.
	@param {!module:private/component/app/root~Address} cell - The address
	of the cell, relative to the top left one.
	@returns {Undefined}
*/
function fillCell(context, cell) {
	const coordinate = cellToCoordinate.call(this, cell);
	const bottomY = coordinate.y + this.scaleHeight;

	if (cellIsHyrule(cell)) {
		context.moveTo(coordinate.x, coordinate.y);
		context.lineTo(coordinate.x - this.scaleWidth, bottomY);
		context.lineTo(coordinate.x + this.scaleWidth, bottomY);
	} else {
		context.moveTo(coordinate.x, bottomY);
		context.lineTo(coordinate.x - this.scaleWidth, coordinate.y);
		context.lineTo(coordinate.x + this.scaleWidth, coordinate.y);
	}

	context.closePath();
	context.fill();
	context.beginPath(); // resetting; see the bottom comment
}

/**
	fadeInCell lets a cell fade in. This function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@param {!external:DOM~CanvasRenderingContext2D} context - The context
	where the cell will be filled.
	@param {!module:private/component/app/root~Address} cell - The address
	of the cell to be filled, relative to the top left one.
	@returns {Undefined}
*/
function fadeInCell(context, cell) {
	const start = Date.now();

	const draw = () => {
		const opacityMax = 255;
		const opacity = Math.min(Date.now() - start, opacityMax);
		const opacityStr = opacity.toString(16);

		context.fillStyle = `#${"0".repeat(8 - opacityStr.length)}${opacityStr}`;
		fillCell.call(this, context, cell);

		if (opacity < opacityMax) {
			requestAnimationFrame(draw);
		}
	};

	requestAnimationFrame(draw);
}

/**
	drawStaticAureole draws static aureole. This function is usable after
	init.
	@private
	@this module:private/component/app/root~Context
	@param {!external:DOM~CanvasRenderingContext2D} context - The context
	where the aureole will be rendered.
	@returns {Undefined}
*/
function drawStaticAureole(context) {
	const middleX = this.width / 2;
	const middleY = this.cellOrigin.y + this.scaleHeight * 4.5;
	const gradient = context.createRadialGradient(
		middleX, middleY, this.scaleWidth * 8,
		middleX, middleY, 0);

	gradient.addColorStop(0, "#000");
	gradient.addColorStop(1, "#880");

	context.fillStyle = gradient;
	context.fillRect(0, 0, this.width, this.height);
}

/**
	newAureoleLines returns a set of new AureoleLine.
	@private
	@returns {!module:private/component/app/root~AureoleLine[]} A set of
	new AureoleLine.
*/
function newAureoleLines() {
	const lines = Array(32);
	const now = Date.now();

	for (let index = 0; index < lines.length; index++) {
		lines[index] = {
			date: now - Math.random() * lineElapsedToScale * lineScale * 2,
			rad:  Math.random() * Math.PI * 2,
		};
	}

	return lines;
}

/**
	newAureoleLineGradient returns a gradient for lines of aureole. This
	function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@param {!external:DOM~CanvasRenderingContext2D} context - The context
	where the gradient creates.
	@returns {!external:DOM~CanvasGradient} A gradient for lines of aureole.
*/
function newAureoleLineGradient(context) {
	const middleX = this.width / 2;
	const middleY = this.cellOrigin.y + this.scaleHeight * 4.5;
	const gradient = context.createRadialGradient(
		middleX, middleY, this.scaleWidth * lineScale,
		middleX, middleY, 0);

	gradient.addColorStop(0, "#88884400");
	gradient.addColorStop(0.8, "#884");
	gradient.addColorStop(0.9, "#88884400");

	return gradient;
}

/**
	showAureole shows aureole. This function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@param {!external:DOM~CanvasRenderingContext2D} context - The context
	where the aureole will be shown.
	@returns {Undefined}
*/
function showAureole(context) {
	drawStaticAureole.call(this, context);

	const background = context.getImageData(0, 0, this.width, this.height);
	const lines = newAureoleLines();
	let aborting = false;
	let opacityPath = 0;
	let opacityStart = 0;
	let opacitySource = 0;

	const draw = () => {
		if (aborting) {
			return;
		}

		const lineLen = this.scaleWidth * lineScale;
		const now = Date.now();
		const opacityProgress = (now - opacityStart) / 16;
		const opacity = opacitySource + Math.sign(opacityPath) * opacityProgress;
		const opacityStr = Math.floor(opacity).toString(16);

		context.putImageData(background, 0, 0);

		for (const line of lines) {
			const lineScaled = (now - line.date) * this.scaleWidth / lineElapsedToScale;
			const middleX = this.width / 2;
			const middleY = this.cellOrigin.y + this.scaleHeight * 4.5;
			let begin;
			let end;

			if (lineScaled < lineLen) {
				begin = 0;
				end = lineScaled;
			} else if (lineScaled < lineLen * 2) {
				begin = lineScaled - lineLen;
				end = lineLen;
			} else {
				line.date = now;
				line.rad = Math.random() * Math.PI * 2;

				continue;
			}

			const sin = Math.sin(line.rad);
			const cos = Math.cos(line.rad);

			context.moveTo(middleX + cos * begin, middleY + sin * begin);
			context.lineTo(middleX + cos * end, middleY + sin * end);
			context.stroke();
			context.beginPath(); // resetting; see the bottom comment
		}

		context.fillStyle = `#${"0".repeat(8 - opacityStr.length)}${opacityStr}`;
		context.fillRect(0, 0, this.width, this.height);

		if (opacityProgress > Math.abs(opacityPath)) {
			opacityPath = Math.random() * 32 - opacity;
			opacityStart = now;
			opacitySource = opacity;
		}

		requestAnimationFrame(draw);
	};

	context.lineWidth = this.scaleHeight / 8;
	context.strokeStyle = newAureoleLineGradient.call(this, context);
	this.abort = () => aborting = true;

	requestAnimationFrame(draw);
}

/**
	drawTriforce draws the Triforce. This function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@param {!external:DOM~CanvasRenderingContext2D} context - The context
	where the Triforce will be drawn.
	@returns {Undefined}
*/
function drawTriforce(context) {
	context.clearRect(0, 0, this.width, this.height);
	context.fillStyle = "#ff0";

	for (const force of triforce) {
		fillCell.call(this, context, force);
	}
}

/**
	playTriforceUnificationSound plays the sound when the Triforce gets
	unified.
	@private
	@returns {Undefined}
*/
function playTriforceUnificationSound() {
	const context = new AudioContext;

	for (const sound of [
		{frequency: 160, gain: 0.04}, {frequency: 320, gain: 0.08},
		{frequency: 525, gain: 0.1}, {frequency: 770, gain: 0.02},
	]) {
		const oscillator = context.createOscillator();
		const gain = context.createGain();

		oscillator.frequency.value = sound.frequency;
		oscillator.connect(gain);

		gain.gain.setValueAtTime(0, context.currentTime);
		gain.gain.linearRampToValueAtTime(sound.gain, context.currentTime + 1);
		gain.gain.linearRampToValueAtTime(0, context.currentTime + 4);
		gain.connect(context.destination);

		oscillator.start(0);
	}

	setTimeout(context.close.bind(context), 4000);
}

/**
	fadeInTriforce fades in the Triforce. This function is usable after the
	trial.
	@private
	@this module:private/component/app/root~Context
	@returns {Undefined}
*/
function fadeInTriforce() {
	const end = () => {
		this.container.children[0].removeEventListener("animationend", end);
		this.container.children[0].style = "position: absolute;";
	};

	this.contexts[0].fillStyle = "#000";
	this.contexts[0].globalCompositeOperation = "destination-over";
	this.contexts[0].fillRect(0, 0, this.width, this.height);

	showAureole.call(this, this.contexts[1]);
	drawTriforce.call(this, this.contexts[2]);

	this.container.children[0].addEventListener("animationend", end);
	this.container.children[0].style.animation = "1s ease 0s 1 normal forwards component-app-root-fading";

	playTriforceUnificationSound();
}

/**
	endTrial ends the trial. This function is usable after the trial.
	@private
	@this module:private/component/app/root~Context
	@returns {Undefined}
*/
function endTrial() {
	const start = Date.now();
	let aborting = false;

	const draw = () => {
		if (aborting) {
			return;
		}

		const maxValue = 0x88; // Why the hell is a supermarket here?
		const value = Math.max(maxValue - (Date.now() - start), 0);

		this.contexts[2].fillStyle = "#" + Math.floor(value).toString(16).repeat(3);
		this.contexts[2].fillRect(0, 0, this.width, this.height);

		this.contexts[0].fillStyle = `hsl(196, ${value / maxValue * 84}%, 56%)`;
		for (const force of triforce) {
			fillCell.call(this, this.contexts[0], force);
		}

		if (value > 0) {
			requestAnimationFrame(draw);
		} else {
			this.abort = null;
			this.contexts[2].globalCompositeOperation = "source-over";
			fadeInTriforce.call(this);
		}
	};

	this.abort = () => aborting = true;
	this.contexts[2].globalCompositeOperation = "darken";
	requestAnimationFrame(draw);
}

/**
	moveUnfixed moves a unfixed cell. This function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@param {!Uint32Array} cells - A bitmap of the cells. See the explanation
	at the bottom of this file
	@param {!module:private/component/app/root~Move} move - A move of the
	cell.
	@returns {Undefined}
*/
function moveUnfixed(cells, move, onmoved) {
	const start = Date.now();

	const draw = () => {
		const [context] = this.contexts;
		const rad = Math.min((Date.now() - start) / 32, Math.PI);
		const middleX = -move.rotation.originRelative.x;
		const topY = -move.rotation.originRelative.y;
		const bottomY = topY + this.scaleHeight;
		let sideY;
		let vertexY;

		if (cellIsHyrule(move.source)) {
			vertexY = topY;
			sideY = bottomY;
		} else {
			vertexY = bottomY;
			sideY = topY;
		}

		context.globalCompositeOperation = "destination-out";

		fillCell.call(this, context, move.destination);
		fillCell.call(this, context, move.source);

		for (const transit of move.transits) {
			fillCell.call(this, context, transit);
		}

		context.fillStyle = "#3be";
		context.globalCompositeOperation = "source-over";

		context.translate(move.rotation.origin.x, move.rotation.origin.y);
		context.rotate(move.rotation.sign * rad / 3);
		context.moveTo(middleX, vertexY);
		context.lineTo(middleX - this.scaleWidth, sideY);
		context.lineTo(middleX + this.scaleWidth, sideY);
		context.closePath();

		context.fill();

		context.beginPath(); // resetting; see the bottom comment
		context.setTransform(1, 0, 0, 1, 0, 0);

		if (rad < Math.PI) {
			requestAnimationFrame(draw);
		} else {
			cells[move.destination.row] |= 0b01 << move.destination.col * 2;
			cells[move.source.row] |= 0b10 << move.source.col * 2;

			for (const transit of move.transits) {
				cells[transit.row] |= 0b10 << transit.col * 2;
			}

			if (onmoved) {
				onmoved();
			}
		}
	};

	cells[move.destination.row] ^= 0b10 << move.destination.col * 2;
	cells[move.source.row] ^= 0b01 << move.source.col * 2;

	for (const transit of move.transits) {
		cells[transit.row] ^= 0b10 << transit.col * 2;
	}

	requestAnimationFrame(draw);
}

/**
	determineMove determines a move of a cell according to the expected
	move of a point in the cell. This function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@param {!module:private/component/app/root~Coordinate} source - The
	absolute coordinate of the point in the cell.
	@param {!module:private/component/app/root~Coordinate} target - The
	absolute coordinate of the target point.
	@returns {!module:private/component/app/root~Move} The move of the
	cell.
*/
function determineMove(source, target) {
	const sourceCoordinate = cellToCoordinate.call(this, source);
	const sourceIsHyrule = cellIsHyrule(source);
	const relativeX = target.x - sourceCoordinate.x;
	const relativeXSign = Math.sign(relativeX);
	const relativeYBySqrt3 = (target.y - sourceCoordinate.y) * sqrt3;
	const transposedYBySqrt3 = sourceIsHyrule ?
		relativeYBySqrt3 - this.scaleWidth * 2 :
		-relativeYBySqrt3 + this.scaleWidth;
	let destination;
	let rotation;
	let transits;

	if (Math.abs(relativeX) < Math.abs(transposedYBySqrt3)) {
		if (relativeX < transposedYBySqrt3) {
			destination = {col: 0, row: 1};

			transits = [
				{col: relativeXSign, row: 0},
				{col: relativeXSign, row: 1},
			];

			rotation = {
				originRelative: {
					x: -relativeXSign * this.scaleWidth,
					y: this.scaleHeight,
				},

				sign: relativeXSign,
			};
		} else {
			destination = {col: relativeXSign, row: 0};

			transits = [
				{col: -relativeXSign, row: 0},
				{col: relativeXSign, row: -1},
			];

			rotation = {
				originRelative: {
					x: relativeXSign * this.scaleWidth,
					y: this.scaleHeight,
				},

				sign: relativeXSign,
			};
		}
	} else {
		destination = {col: relativeXSign, row: 0};
		transits = [{col: 0, row: 1}, {col: relativeXSign * 2, row: 0}];
		rotation = {originRelative: {x: 0, y: 0}, sign: -relativeXSign};
	}

	if (!sourceIsHyrule) {
		destination.row = -destination.row;
		rotation.originRelative.y = this.scaleHeight - rotation.originRelative.y;
		rotation.sign = -rotation.sign;

		for (const transit of transits) {
			transit.row = -transit.row;
		}
	}

	destination.col += source.col;
	destination.row += source.row;

	for (const transit of transits) {
		transit.col += source.col;
		transit.row += source.row;
	}

	rotation.origin = {
		x: sourceCoordinate.x + rotation.originRelative.x,
		y: sourceCoordinate.y + rotation.originRelative.y,
	};

	return {destination, source, transits, rotation};
}

/**
	trialNewCells returns a new bitmap of the cells in a trial.
	@private
	@returns {!Uint32Array} A new bitmap of the cells in the trial. See the
	explanation at the bottom of this file
*/
function trialNewCells() {
	const cells = new Uint32Array(8);

	let margin = 0b101010101010101010101010101010;
	for (let row = 4; row >= 0; row--) {
		cells[row] = margin;
		cells[9 - row] = margin;

		margin = (margin >> 2) & (margin << 2);
	}

	for (const fixed of fixeds) {
		cells[fixed.row] ^= 0b10 << (fixed.col * 2);
	}

	return cells;
}

/**
	startTrial starts a trial. This function is usable after fading in.
	@private
	@this module:private/component/app/root~Context
	@returns {Undefined}
*/
function startTrial() {
	const cells = trialNewCells();
	let trackingCell;
	let startTracking;
	let endTracking;

	const track = event => {
		const eventCell = coordinateToCell.call(this,
		{x: event.layerX, y: event.layerY});

		if (eventCell.col == trackingCell.col && eventCell.row == trackingCell.row) {
			return;
		}

		const move = determineMove.call(this, trackingCell,
			{x: event.layerX, y: event.layerY});

		for (const transit of move.transits) {
			if (!(cells[transit.row] >> transit.col * 2 & 0b10)) {
				return;
			}
		}

		if (!(cells[move.destination.row] >> move.destination.col * 2 & 0b10)) {
			return;
		}

		moveUnfixed.call(this, cells, move, () => {
			if (triforce.every(force => cells[force.row] >> force.col * 2 & 0b01)) {
				this.container.removeEventListener("mousedown", startTracking);
				this.container.removeEventListener("mouseleave", endTracking);
				this.container.removeEventListener("mouseup", endTracking);

				endTrial.call(this);
			}
		});

		endTracking();
	};

	startTracking = event => {
		const cell = coordinateToCell.call(this,
			{x: event.layerX, y: event.layerY});

		if (cells[cell.row] >> cell.col * 2 & 0b01) {
			trackingCell = cell;
			this.container.addEventListener("mousemove", track);
		}
	};

	endTracking = () => this.container.removeEventListener("mousemove", track);

	this.contexts[0].clearRect(0, 0, this.width, this.height);
	this.contexts[0].fillStyle = "#3be";

	for (const unfixed of initUnfixeds) {
		fillCell.call(this, this.contexts[0], unfixed);
		cells[unfixed.row] ^= 0b11 << (unfixed.col * 2);
	}

	this.container.children[0].style.zIndex = "1";
	this.container.addEventListener("mousedown", startTracking);
	this.container.addEventListener("mouseleave", endTracking);
	this.container.addEventListener("mouseup", endTracking);
}

/**
	fadeInTrial fades in the trial. This function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@returns {Undefined}
*/
function fadeInTrial() {
	const start = Date.now();
	let aborting = false;

	const draw = () => {
		if (aborting) {
			return;
		}

		const context = this.contexts[2];
		const {cellOrigin} = this;
		const alpha = (Date.now() - start) * 0.004;
		const topY = cellOrigin.y;
		const middleY = cellOrigin.y + this.scaleHeight * 5;
		const bottomY = cellOrigin.y + this.scaleHeight * 8;

		context.globalAlpha = alpha;

		context.fillStyle = "#000";
		context.fillRect(0, 0, this.width, this.height);

		context.fillStyle = "#888";
		context.moveTo(cellOrigin.x + this.scaleWidth * 4, topY);
		context.lineTo(cellOrigin.x + this.scaleWidth * 10, topY);
		context.lineTo(cellOrigin.x + this.scaleWidth * 15, middleY);
		context.lineTo(cellOrigin.x + this.scaleWidth * 12, bottomY);
		context.lineTo(cellOrigin.x + this.scaleWidth * 2, bottomY);
		context.lineTo(cellOrigin.x - this.scaleWidth, middleY);
		context.closePath();
		context.fill();
		context.beginPath(); // resetting; see the bottom comment

		context.fillStyle = "#adf";
		for (const force of triforce) {
			fillCell.call(this, context, force);
		}

		context.globalAlpha = 1;
		context.fillStyle = "#000";
		for (const fixed of fixeds) {
			fillCell.call(this, context, fixed);
		}

		if (alpha < 1) {
			requestAnimationFrame(draw);
		} else {
			this.abort = null;
			startTrial.call(this);
		}
	};

	this.abort = () => aborting = true;
	this.contexts[2].globalAlpha = 0;
	this.contexts[2].globalCompositeOperation = "source-over";

	requestAnimationFrame(draw);
}

/**
	swap swaps slides in canvas and calls fadeInTrial accordingly. This
	function is usable after init.
	@private
	@this module:private/component/app/root~Context
	@returns {Undefined}
*/
function swap() {
	let {x} = this.tileOrigin;
	let background = 0;
	let foreground = 1;
	let slide = 0;
	let fadingInFixeds;
	let interval;
	let start;

	const lineEdge = () => {
		const context = this.contexts[foreground];

		context.moveTo(x - this.scaleWidth, this.tileOrigin.y);

		let cursorX = x + this.scaleWidth * 2;
		let cursorY = this.tileOrigin.y;
		while (cursorX >= 0 && cursorY <= this.width) {
			context.lineTo(cursorX - this.scaleWidth, cursorY);
			cursorY += this.scaleHeight;
			context.lineTo(cursorX, cursorY);
			cursorX -= this.scaleWidth;
		}

		context.lineTo(cursorX - this.scaleWidth * 3, cursorY);
	};

	const lineDoubleEdge = (initX, initY) => {
		const context = this.contexts[foreground];
		let cursorX = initX;
		let cursorY = initY;

		while (true) {
			context.lineTo(cursorX, cursorY);
			cursorX -= this.scaleWidth * 2;

			if (cursorX <= 0 || cursorY + this.scaleHeight * 2 > this.width) {
				break;
			}

			context.lineTo(cursorX - this.scaleWidth * 2, cursorY);
			cursorY += this.scaleHeight * 2;
		}

		context.lineTo(cursorX, cursorY);
	};

	const lineDoubleEdge0 = () => {
		const context = this.contexts[foreground];

		context.moveTo(x + this.scaleWidth, this.tileOrigin.y);

		lineDoubleEdge(x + this.scaleWidth * 2,
			this.tileOrigin.y + this.scaleHeight);
	};

	const lineDoubleEdge1 = () => {
		const context = this.contexts[foreground];

		context.moveTo(x - this.scaleWidth, this.tileOrigin.y);

		lineDoubleEdge(x + this.scaleWidth,
			this.tileOrigin.y + this.scaleHeight * 2);
	};

	const restart = () => {
		slide++;
		if (slide >= this.slideImages.length) {
			slide = 0;
		}

		const oldForeground = foreground;

		foreground = background;
		background = oldForeground;

		this.container.children[oldForeground].style = "position: absolute;";
		this.container.children[foreground].style.zIndex = "1";

		setTimeout(start, 8192);
	};

	const stop = () => clearInterval(interval);

	const click = event => {
		if (event.layerX < x - event.layerY / sqrt3) {
			const cell = coordinateToCell.call(this,
				{x: event.layerX, y: event.layerY});

			fixeds.some((fixed, index) => {
/* FOR DEBUGGING (and cheating :P)
				fadingInFixeds |= 1 << index;

				fadeInCell.call(this,
					this.contexts[foreground], fixed);
*/
				if (cell.col == fixed.col && cell.row == fixed.row) {
					const bit = 1 << index;

					if (!(fadingInFixeds & bit)) {
						fadingInFixeds |= bit;

						fadeInCell.call(this,
							this.contexts[foreground],
							fixed);
					}

					return true;
				}
			});
		}
	};

	const draw = () => {
		const context = this.contexts[foreground];
		const foregroundDOM = this.container.children[foreground];
		const backgroundDOM = this.container.children[background];

		context.globalCompositeOperation = "destination-out";
		lineEdge();
		context.closePath();
		context.fill();
		context.globalCompositeOperation = "source-over";

		lineDoubleEdge0();
		context.closePath();
		lineDoubleEdge1();
		context.stroke();
		context.beginPath(); // resetting; see the bottom comment

		x += this.scaleWidth * 2;
		if (x - this.height / sqrt3 >= this.width) {
			stop();

			if (fadingInFixeds == (1 << fixeds.length) - 1) {
				for (const dom of [backgroundDOM, foregroundDOM]) {
					dom.removeEventListener("animationend",
						restart);
				}

				foregroundDOM.style.zIndex = "";
				this.container.removeEventListener("click", click);
				this.stop = null;
				fadeInTrial.call(this);
			} else {
				foregroundDOM.style.animation = "component-app-root-fading 1s";
				x = this.tileOrigin.x;
			}
		}
	};

	start = () => {
		const slideRatio = this.slideImages[slide].width / this.slideImages[slide].height;
		let width;
		let height;

		fadingInFixeds = 0;

		if (slideRatio < this.width / this.height) {
			width = this.width;
			height = this.width / slideRatio;
		} else {
			width = this.height * slideRatio;
			height = this.height;
		}

		this.contexts[background].drawImage(this.slideImages[slide],
			0, 0, width, height);

		interval = setInterval(draw, 256);
	};

	for (const index of [background, foreground]) {
		this.container.children[index].addEventListener("animationend", restart);
		this.contexts[index].strokeStyle = "#fff";
	}

	this.abort = stop;
	this.contexts[foreground].fillStyle = "#fff";
	this.contexts[foreground].fillRect(0, 0, this.width, this.height);

	this.container.addEventListener("click", click);
	start();
}

/**
	init initializes. This function is usable when container and swapImages
	are set. It does initialize the following property and starts swapping:
	contexts, cellOrigin, tileOrigin, scaleWidth, scaleHeight, width, and
	height.
	@private
	@this module:private/component/app/root~Context
	@returns {Undefined}
*/
function init() {
	const {children} = this.container;
	const jQuery = $(this.container);

	this.contexts = Array.prototype.map.call(children,
		child => child.getContext("2d"));

	this.width = jQuery.width();
	this.height = jQuery.height();

	const widthWidth = this.width / 16;
	const widthHeight = this.width * sqrt3;
	const heightHeight = this.height / 8;

	if (widthHeight < heightHeight) {
		const middleY = this.height / 2;

		this.scaleWidth = widthWidth;
		this.scaleHeight = widthHeight;

		this.cellOrigin = {
			x: this.scaleWidth,
			y: middleY - this.scaleHeight * 4,
		};

		this.tileOrigin = {
			x: 0,
			y: middleY % this.scaleHeight - this.scaleHeight,
		};
	} else {
		this.scaleWidth = heightHeight / sqrt3;
		this.scaleHeight = heightHeight;

		const cellWidth = this.scaleWidth * 2;
		const middleX = this.width / 2;

		this.cellOrigin = {
			x: middleX - this.scaleWidth * 7,
			y: 0,
		};

		this.tileOrigin = {
			x: middleX % cellWidth - cellWidth,
			y: 0,
		};
	}

	for (const child of children) {
		child.style.position = "absolute";
		child.width = this.width;
		child.height = this.height;
	}

	swap.call(this);
}

/**
	setContainer sets container property in Context and calls init if it is
	ready.
	@private
	@this module:private/component/app/root~Context
	@param {!external:Mithril~Node} - A virtual DOM node of the container.
	@returns {Undefined}
*/
function setContainer(node) {
	this.container = node.dom;

	if (this.slideImages.length) {
		init.call(this);
	}
}

export function oninit() {
	this.slideImages = [];

	slides.forEach(slide => {
		const slideImage = new Image;
		slideImage.src = slide;

		slideImage.addEventListener("load", () => {
			this.slideImages.push(slideImage);

			if (this.slideImages.length == 1 && this.container) {
				init.call(this);
			}
		});
	});
}

export function onbeforeremove() {
	if (this.abort) {
		this.abort();
	}
}

export function view() {
	return m(container, m("div", {
		oncreate: setContainer.bind(this),
		style:    {height: "100%", width: "100%"},
	}, m("canvas"), m("canvas"), m("canvas")));
}

/**
	Address is an address of a cell.
	@private
	@typedef module:private/component/app/root~Address
	@property {!Number} col - The column of the cell, relative to the left.
	@property {!Number} row - The row of the cell, relative to the top.
*/

/**
	Coordinate is a coordinate in a canvas.
	@private
	@typedef module:private/component/app/root~Coordinate
	@property {!Number} x - The path from the vertical bar as the reference.
	@property {!Number} y - The path from the horizonal bar as the
	reference.
*/

/**
	AureoleLine is a line in aureole.
	@private
	@typedef module:private/component/app/root~AureoleLine
	@property {!Number} rad - The radian.
	@property {!Number} date - The date when the AureoleLine got
	initialized.
*/

/**
	Move is a move of a cell.
	@private
	@typedef module:private/component/app/root~Move
	@property {!module:private/component/app/root~Rotation} rotation - The
	rotation.
	@property {!module:private/component/app/root~Address} destination -
	The address of the destination, relative to the top left cell.
	@property {!module:private/component/app/root~Address} source -
	The address of the source, relative to the top left cell.
	@property {!module:private/component/app/root~Address[]} transits -
	A set of the addresses of the transits.
*/

/**
	Rotation is a rotation of a cell.
	@private
	@typedef module:private/component/app/root~Rotation
	@property {!module:private/component/app/root~Coordinate} origin - The
	absolute coordinate to the center of the rotation.
	@property {!module:private/component/app/root~Coordinate}
	originRelative - The coordinate to the center of the rotation, relative
	to the coordinate of the source cell.
	@property {!Number} sign - The sign of the radian determining the
	direction of the rotation.
*/

/**
	Context is a context persistent in the component. TODO
	@private
	@typedef module:private/component/app/root~Context
	@property {?external:DOM~HTMLElement} container - A container of the
	canvases.
	@property {?function} abort - A function to abort the rendering. After
	calling this function, the relevant DOM resources can safely be
	released. If it is null, the relevant DOM resources can safely be
	released without any preparation.
	@property {?external:DOM~Canvas2DContext[]} contexts - Canvas2DContexts.
	Each contexts corresponds to the element of container#children with the
	same index.
	@property {?module:private/component/app/root~Coordinate} cellOrigin -
	The absolute coordinate of the left top cell.
	@property {?module:private/component/app/root~Coordinate} tileOrigin -
	The absolute left top coordinate of the lined box. That could be
	negative in order to include extra parts of cells.
	@property {!external.DOM~HTMLImageElement[]} slideImages - Loaded slide
	images.
	@property {?Number} scaleWidth - The factor for scaling the width of
	the cells. It is the half number of the cells' width.
	@property {?Number} scaleHeight - The factor for scaling the height of
	the cells. It is actually same with the cells' height.
	@property {?Number} width - The width of the contents of the canvases
	and container.
	@property {?Number} height - The height of the contents of the canvases
	and container.
*/

/**
	Trial is a context of a trial.
	@private
	@typedef module:private/component/app/root~Trial
*/

/*
	An explanation of a bitmap of the cells (referred as "the explanation at
	the bottom of this file"):
	A bitmap of the cells is a Uint32Array describing the type of the cells.
	Each elements represents a 16-cells-wide row.
	In an element, 2-bits wide cell descriptors are stored without any
	margin between them.
	A descriptor in lesser significant bits represents cells more close to
	the left side.
	The descriptor could be 0b00, 0b01, or 0b10, which repectively means
	a fixed cell, an unfixed cell, or nothing.

	Notice about beginPath (referred as "the bottom comment"):
	beginPath sounds as if that should be done before lining. However, that
	is FALSE. The standard says:

	HTML Standard 4.12.5.1.12 Drawing paths to the canvas
	https://html.spec.whatwg.org/multipage/scripting.html#dom-context-2d-beginpath
	> The beginPath() method, when invoked, must empty the list of subpaths
	> in the context's current default path so that the it once again has
	> zero subpaths.

	Here, we have three options about dealing with the state:
	1. Initialize the state anytime after modifying it.
	2. Initialize the state anytime before modifying it.
	3. Do both.

	2 and 3 are not realistic since Canvas2DContext carries lots of
	variables, and, importantly the number could even be increase in the
	future. Moreover, persisting states may make extra memory usage.
	We choose 1; take care of your sh*ts by yourself!
*/

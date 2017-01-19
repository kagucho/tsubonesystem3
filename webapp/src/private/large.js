export default function() {
	const container = $("#container");

	return container.width() / parseFloat(container.css("font-size")) > 64;
}

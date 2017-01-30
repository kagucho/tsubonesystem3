export function view() {
	return m("div", {
		className: "modal fade",
		config(element, initialized) {
			const jquery = $(element);

			if (!initialized) {
				jquery.on("shown.bs.modal", () => setTimeout(() => console.log(m.redraw()), 1000));
				jquery.modal("show");
			}
		},
		role:     "dialog",
		tabindex: "-1",
	}, m("div", {className: "modal-dialog", role: "document"},
			m("form", {className: "modal-content"},
				m("div", {className: "modal-body"},
					"Hello, World!"
				), m("div", {className: "modal-footer"},
					m("button", {
						className: "btn btn-default",
						onclick:   m.redraw,
					}, "REDRAW")
				)
			)
		)
	);
}

module.exports = {
	plugins: [
		require("stylelint")({
			config: {
				extends: "stylelint-config-standard",
				rules:   {
					"indentation":                       "tab",
					"number-leading-zero":               "never",
					"selector-list-comma-newline-after": "never-multi-line",
				},
			},
		}),
		require("postcss-cssnext")({
			browsers: [
				"last 1 Chrome versions",
				"Edge >= 20",
				"Firefox ESR",
				"ie >= 11",
			],
		}),
	],
};

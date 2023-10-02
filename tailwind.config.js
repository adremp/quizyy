const plugin = require("tailwindcss/plugin");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./node_modules/tw-elements/dist/js/**/*.js"],
  theme: {
    extend: {
      keyframes: {
        show: {
          "0%": {
            opacity: "0",
          },
          "100%": {
            opacity: "1",
          },
        },
        "move-left-right": {
          "0%, 100%": { transform: "translateX(0)" },
          "25%": { transform: "translateX(-25%)" },
          "75%": { transform: "translateX(25%)" },
        },
      },
      animation: {
        show: "show 0.3s ease",
        "move-left-right": "move-left-right 0.4s ease",
      },
    },
  },
  plugins: [
    plugin(({ addVariant, matchVariant, addComponents }) => {
      addVariant("child", "& > *");
      addComponents({
        ".disabled": {
          opacity: 0.5,
          pointerEvents: "none",
        },
        ".plus": {
          position: "relative",
          display: "inline-flex",
					alignItems: "center",
					justifyContent: "center",
          height: "100%",
          aspectRatio: "1",
          "&::before, &::after": {
            content: "",
            position: "absolute",
            backgroundColor: "currentColor",
            width: "100%",
            height: "100%",
          },
					"&::before": {
						rotate: "45deg"
					},
					"&::after": {
						rotate: "-45deg"
					}
        },
      });
      matchVariant("has", (val) => `&:has(${val})`);
      matchVariant("group-has", (val) => `:merge(.group):has(${val}) &`);
    }),
    require("tw-elements/dist/plugin.cjs"),
  ],
};

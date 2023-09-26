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
    plugin(({ addVariant }) => {
      addVariant("child", "& > *");
    }),
    require("tw-elements/dist/plugin.cjs"),
  ],
};


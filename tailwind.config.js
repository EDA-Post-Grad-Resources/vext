/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [
    function ({ addUtilities }) {
      addUtilities({
        ".diagonal-lines": {
          background:
            "repeating-linear-gradient(45deg, rgba(0,0,0,1) 1px, rgba(0,60,0,1) 2px )",
        },
      });
    },
  ],
};

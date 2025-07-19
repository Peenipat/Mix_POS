module.exports = {
    content: [
      "./index.html",
      "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
      extend: {
        colors: {
          primary: {
            light: "#3B82F6",
            DEFAULT: "#2563EB",
            dark: "#1D4ED8",
          },
        },
      },
    },
    darkMode: "class", 
    plugins: [
      require('flowbite/plugin')
    ],
  };
  
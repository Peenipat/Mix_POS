module.exports = {
  content: [
    './src/**/*.{js,ts,jsx,tsx}',
    './node_modules/flowbite-react/**/*.js',
    './node_modules/flowbite/**/*.js',
    './node_modules/flowbite-datepicker/**/*.js', 
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"K2D"', 'sans-serif'], 
      },
    },
  },
  plugins: [
    require('flowbite/plugin'),
    require('tailwindcss-filters'),
    require('@tailwindcss/forms'),
    require('flowbite-datepicker/plugin'),
  ],
};

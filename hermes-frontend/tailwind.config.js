/** @type {import('tailwindcss').config} */
module.exports = {
  important: true,
  content: [
    "./src/**/*.{js,jsx, ts,tsx}",
  ],
  theme: {
    container: false,
    screens: {
      'tablet': {'max': '1180px'},
      'phablet': {'max': '980px'},
      'phone': {'max': '670px'},
      'tiny': {'max': '450px'},
    },
    colors: {
      'background': '#f4faff',
      'background-secondary': '#fff',
      'dark': '#191919',
      'dark-secondary': '#535657',
      'dark-highlight': '#4f656f',
      'accent': '#4583f6',
      'primary': '#345995',
      'error': '#f45866',
      'error-accent': '#f56860',
      'success': '#44C26E',
      'success-accent': '#44C26E',
    },
    fontFamily: {
      sans: ['karla', 'sans-serif']
    },
    fontSize: {
      sm: '0.8rem',
      base: '18px',
      lg: '1.25rem',
      'xl': '1.563rem',
      '2xl': '1.953rem',
      '3xl': '2.441rem',
      phablet: '16px',
      phone: '14px',
    },
    extend: {
      borderRadius: {
        'br': '8px'
      },
      boxShadow: {
        'bs': '#e7e7e7 1px 1px 5px',
        'emp': '#d4d4d4 1px 1px 7px',
      },
      dropShadow: {
        'ds': '#e7e7e7 1px 1px 5px',
      }
    },
  },
  plugins: [
  ],
}

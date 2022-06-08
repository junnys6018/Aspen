const defaultTheme = require('tailwindcss/defaultTheme');

module.exports = {
    content: ['./pages/**/*.{js,ts,jsx,tsx}', './components/**/*.{js,ts,jsx,tsx}'],
    theme: {
        fontFamily: {
            sans: ['Poppins', 'system-ui'],
            mono: defaultTheme.fontFamily.mono,
        },
        container: {
            padding: {
                DEFAULT: '1rem',
                sm: '2rem',
                lg: '2rem',
                xl: '10rem',
                '2xl': '16rem',
            },
            center: true,
        },
        extend: {
            spacing: {
                160: '40rem',
            },
        },
    },
    plugins: [require('@tailwindcss/typography')],
};

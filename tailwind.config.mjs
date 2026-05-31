/** @type {import('tailwindcss').Config} */
/** this is a v3 theme config that exists mostly to bamboozle the LSP into giving me classname hints defined in inlined CSS in a go html template */

module.exports = {
	theme: {
		extend: {
			screens: {
				xs: '26rem'
			},
			boxShadow: {
				'up': '0 -10px 15px -3px rgba(0, 0, 0, 0.1)',
			},
			colors: {
				'tfl-blue': 'rgb(30,41,59)',
				'tfl-gray': 'rgb(71,85,105)',
				'whatsapp-green': '#25D366',
				'whatsapp-green-dark': '#20ba5a',
			},
			animation: {
				'fade-in-up': 'fadeInUp 0.2s ease-out forwards',
			},
		}
	}
}

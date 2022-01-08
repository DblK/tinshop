import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://github.com/sveltejs/svelte-preprocess
	// for more information about preprocessors
	preprocess: preprocess(),

	kit: {
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			// fallback: null,
			// precompress: false,
			// fallback: '200.html',
			fallback: 'index.html',
		}),
		paths: {
			base: '/admin'
		},
		amp: false,
		appDir: 'internal',
		ssr: false,
		browser: true,

		// hydrate the <div id="svelte"> element in src/app.html
		target: '#svelte'
	}
};

export default config;

import { sveltekit } from '@sveltejs/kit/vite';
import basicSsl from '@vitejs/plugin-basic-ssl'
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [basicSsl(), sveltekit()],
	server: {
		https: true,
		host: '0.0.0.0',
		proxy: {}
	}
});

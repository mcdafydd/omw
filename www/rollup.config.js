import { createDefaultConfig } from '@open-wc/building-rollup';
import cpy from 'rollup-plugin-cpy';
import litcss from 'rollup-plugin-lit-css';

const config = createDefaultConfig({ input: './index.html' });

export default {
  ...config,
  output: {
    ...config.output,
    sourcemap: false,
  },
  plugins: [
    ...config.plugins,
    litcss(),
    cpy({
      files: ['images/*.png', 'images/*.svg', 'components/styles.css', 'favicon.ico', 'manifest.json',
	      'node_modules/@fullcalendar/core/main.css',
              'node_modules/@fullcalendar/daygrid/main.css',
              'node_modules/@fullcalendar/timegrid/main.css',
      ],
      dest: 'dist',
      options: {
        // parents makes sure to preserve the original folder structure
        parents: true,
      },
    }),
  ],
};

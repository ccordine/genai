module.exports = {
  apps: [
    {
      name: 'genai',
      script: './genai',
      interpreter: 'none',
      exp_backoff_restart_delay: 100,
      autorestart: true,
      watch: false,
      max_memory_restart: '300M',
      env: {
        NODE_ENV: 'production',
      },
    },
  ],
};


// PM2 process file for the BaliGuide API.
//
//   make pm2-start     build, then start (or reload) under PM2
//   make pm2-restart   rebuild after code changes and restart
//   make pm2-logs      follow the logs
//   make pm2-stop      stop the API
//
// PM2 runs the compiled binary (bin/api), not `go run`. The working
// directory is this folder so the API finds .env and the uploads/ folder.

module.exports = {
  apps: [
    {
      name: "baliguide-api",
      script: "./bin/api",
      cwd: __dirname,
      interpreter: "none", // a Go binary, not a Node script
      exec_mode: "fork",
      instances: 1, // rate limits and caches live in memory: keep one process

      autorestart: true,
      max_restarts: 10,
      min_uptime: "10s", // a crash sooner than this counts as a failed start
      restart_delay: 3000,
      kill_timeout: 5000,

      out_file: "./logs/api.out.log",
      error_file: "./logs/api.err.log",
      merge_logs: true,
      time: true, // prefix log lines with a timestamp

      // Settings come from .env (loaded by the API itself). Values here
      // override .env only when set.
      env: {
        APP_ENV: "development",
      },
      env_production: {
        APP_ENV: "production",
        GIN_MODE: "release",
      },
    },
  ],
};

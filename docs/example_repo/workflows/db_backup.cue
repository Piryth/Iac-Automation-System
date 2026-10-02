package ias

workflows: "daily-backup": {
	name: "Database Backup"
	trigger: {
		type: "cron"
		config: schedule: "0 2 * * *"
	}
	graph: {
		entrypoint: "dump_db"
		nodes: {
			"dump_db": {
				type: "command"
				config: command: "pg_dump mydb > /tmp/db.sql"
			}
		}
	}
}
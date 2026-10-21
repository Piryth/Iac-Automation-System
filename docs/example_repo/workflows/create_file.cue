package ias

workflows: "pings": {
	name: "Files ping"
	trigger: {
		type: "cron"
		config: schedule: "* * * * *"
	}
	graph: {
		entrypoint: "create_file"
		nodes: {
			"create_file": {
				type: "command"
				config: command: "touch '/home/bl204483/timestamps/created_$(date +%Y%m%d_%H%M%S)'"
			}
		}
	}
}
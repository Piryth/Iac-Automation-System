package ias

#Project: {
	server: address: string | *"localhost:50051"
	workflows: [WorkflowID=string]: #Workflow & {id: WorkflowID}
}

#Workflow: {
	id:          string
	name:        string
	description: string | *""
	trigger:     #Trigger
	graph:       #Graph
}

#Trigger: {
	type: "cron"
	config: schedule: string
}

#Graph: {
	entrypoint: string
	nodes: [NodeID=string]: #Node & {id: NodeID}
}

#Node: {
	id:         string
	type:       "command"
	on_success?: string
	config: command: string
}

// Root evaluation constraint
#Project
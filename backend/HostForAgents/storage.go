package HostForAgents

import "diplom_server/backend/structs"

type FollowedAgent struct {
	CurrentState structs.Stats
	Changes      structs.Stats
}

type Storage struct {
	Hosts []structs.Host
	Data  []structs.Stats

	FollowedAgents map[string]FollowedAgent
}

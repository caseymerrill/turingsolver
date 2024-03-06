package game

type RemotePlayer struct {
	Name string `json:"name"`
}

func (p *RemotePlayer) GetPlayerName() string {
	return p.Name
}
func (p *RemotePlayer) Solve(g Game) (bool, Solution) {
	return false, Solution{}
}
func (p *RemotePlayer) SetProgressCallback(callback ProgressCallback) {}
func (p *RemotePlayer) Clone() Player {
	return &RemotePlayer{Name: p.Name}
}

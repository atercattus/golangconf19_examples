package scenes

var (
	Scenes struct {
		Title      *TitleScene
		SelectName *SelectNameScene
		PartyMaker *PartyMakerScene
		Game       *GameScene
	}
)

func init() {
	Scenes.Title = &TitleScene{}
	Scenes.Title.Name = `title`

	Scenes.SelectName = &SelectNameScene{}
	Scenes.SelectName.Name = `select_name`

	Scenes.PartyMaker = &PartyMakerScene{}
	Scenes.PartyMaker.Name = `party_maker`

	Scenes.Game = &GameScene{}
	Scenes.Game.Name = `game`
}

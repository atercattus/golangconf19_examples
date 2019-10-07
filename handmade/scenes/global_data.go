package scenes

import (
	"github.com/atercattus/golangconf19_examples/handmade/events"
	"github.com/atercattus/golangconf19_examples/handmade/net"
)

type (
	GlobalDataState struct {
		PlayerName       string
		PlayerId         net.Id
		RaceId           net.Id
		RaceDistance     float32
		Players          []events.PlayerInfoWithTimes
		PlayersOurIdx    int // в какой позиции в Players лежат данные по текущему игроку
		LastStepAtServer net.Time
		LastStepAtPlayer net.Time

		MaxPlayers    int
		pixelsPerStep float32

		WS net.WebsocketClienter
	}
)

const (
	ScreenWidthInSteps = 70
)

var (
	GlobalData GlobalDataState
)

func (gd *GlobalDataState) Init() {
	gd.PlayerId = net.GenId()
}

func (gd *GlobalDataState) SetPixelsPerStep(pps float32) {
	gd.pixelsPerStep = pps
}

func (gd *GlobalDataState) GetPixelsPerStep() float32 {
	return gd.pixelsPerStep
}

func (gd *GlobalDataState) Join(username string, cb net.EventCallback) error {
	gd.PlayerName = username

	var ev events.EventJoin
	ev.PlayerName = username
	ev.PlayerId = gd.PlayerId

	return gd.WS.SendMessage(events.EventToNetEvent(&ev), cb)
}

func (gd *GlobalDataState) SendIAmReady(cb net.EventCallback) error {
	var ev events.EventIAmReady
	ev.PlayerId = gd.PlayerId
	ev.RaceId = gd.RaceId

	return gd.WS.SendMessage(events.EventToNetEvent(&ev), cb)
}

func (gd *GlobalDataState) SendIAmFinished(cb net.EventCallback) error {
	var ev events.EventIAmFinished
	ev.PlayerId = gd.PlayerId
	ev.RaceId = gd.RaceId

	return gd.WS.SendMessage(events.EventToNetEvent(&ev), cb)
}

func (gd *GlobalDataState) playerStep() error {
	player := &gd.Players[gd.PlayersOurIdx]

	gd.LastStepAtPlayer = net.GetCurrentTime()

	var ev events.EventPlayerStep
	ev.RaceId = gd.RaceId
	ev.PlayerId = gd.PlayerId
	ev.Distance = player.Distance
	ev.Speed = player.Speed
	ev.Now = gd.LastStepAtPlayer

	return gd.WS.SendMessage(events.EventToNetEvent(&ev), nil)
}

func (gd *GlobalDataState) SetPlayersList(players []events.PlayerInfo) {
	gd.Players = gd.Players[:0]

	gd.PlayersOurIdx = -1 // всегда должен быть
	for idx, player := range players {
		gd.Players = append(gd.Players, events.PlayerInfoWithTimes{
			PlayerInfo: player,
		})

		if player.Id == gd.PlayerId {
			gd.PlayersOurIdx = idx
		}
	}
}

func (gd *GlobalDataState) UpdatePlayersState(playersStates []events.PlayerCompactState) (ourFinishPos int, smbdFinished bool,
) {
	ourFinishedAt := net.Time(0)

	for i, playerState := range playersStates {
		if playerState.FinishedAt > 0 {
			smbdFinished = true
		}

		if i == gd.PlayersOurIdx {
			gd.LastStepAtServer = playerState.LastStepAtServer
			ourFinishedAt = playerState.FinishedAt
			continue
		}
		player := &gd.Players[i]
		player.Distance = playerState.Distance
		player.Speed = playerState.Speed
		player.LastStepAtServer = playerState.LastStepAtServer
		player.FinishedAt = playerState.FinishedAt
	}
	if err := gd.playerStep(); err != nil {
		println(`Send playerStep error:`, err.Error())
	}

	if ourFinishedAt > 0 {
		ourFinishPos = 1
		for _, playerState := range playersStates {
			if (playerState.FinishedAt > 0) && (playerState.FinishedAt < ourFinishedAt) {
				ourFinishPos++
			}
		}
	}

	return ourFinishPos, smbdFinished
}

// Предсказание, где сейчас находится другой игрок на основе его старых сведений
func (gd *GlobalDataState) PredictPlayerDistance(playerInfo *events.PlayerInfoWithTimes) float32 {
	// Когда сервер в последний раз получал от нас событие EventPlayerStep
	tsOurLastStepServerTime := float64(gd.LastStepAtServer) / 1000 // t3
	// Когда мы в последний раз отсылали серверу событие EventPlayerStep
	tsOurLastStepLocalTime := float64(gd.LastStepAtPlayer) / 1000 // t1
	// Когда сервер в последний раз получал от другого игрока событие EventPlayerStep
	tsOtherLastStepServerTime := float64(playerInfo.LastStepAtServer) / 1000 // t4
	// Текущее локальное время у нас
	tsNowLocalTime := float64(net.GetCurrentTime()) / 1000

	dist := playerInfo.Distance
	if (tsOtherLastStepServerTime == 0) || (tsOurLastStepServerTime == 0) || (tsOurLastStepLocalTime == 0) {
		return dist
	}

	// Когда другой игрок в последний раз отсылал серверу событие EventPlayerStep
	tsOtherLastStepLocalTime := tsOtherLastStepServerTime - tsOurLastStepServerTime + tsOurLastStepLocalTime // t2

	// Сколько примерно времени прошло с момента последнего EventPlayerStep у другого игрока по его версии
	timeElapsed := tsNowLocalTime - tsOtherLastStepLocalTime

	// Мы не знаем, была ли еще какая активность от другого игрока, так что просто считаем как есть
	return playerInfo.Distance + playerInfo.Speed*float32(timeElapsed)
}

func (gd *GlobalDataState) GetTimeDiffWithServerInMsec() int {
	const maxRealDiff = 5000 // 5 секунд
	if GlobalData.LastStepAtPlayer == 0 || GlobalData.LastStepAtServer == 0 {
		return 0
	}
	diff := int(GlobalData.LastStepAtPlayer - GlobalData.LastStepAtServer)
	if diff < -maxRealDiff || diff > maxRealDiff {
		diff = 0
	}
	return diff
}

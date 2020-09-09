package handlers

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// VoteStartHandler handles vote start requests.
type VoteStartHandler struct{}

// RequestType returns Vote Start Request identifier.
func (*VoteStartHandler) RequestType() data.MessageType {
	return data.MessageVoteStartRequest
}

// BodyReader parses VoteStartRequest json from ws.
func (*VoteStartHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.VoteStartRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes vote start requests.
func (*VoteStartHandler) Handle(server *web.Server, req *web.Request) {
	voteStart, ok := req.Body.(*data.VoteStartRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting VoteStartRequest, found", req.Body)
		return
	}

	g, err := server.VoteStart(req.UserIO, voteStart.Vote)

	if err != nil {
		req.SendError(err)

		return
	}

	req.SendMessage(data.MessageEmptyResponse, nil)

	if g == nil {
		return
	}

	turnSequence := g.PlayerTurnSequence()

	server.NotifyPlayers(g, nil, data.MessageNotifyGameStarted, func(player *game.Player) interface{} {
		return data.NotifyGameStarted{
			PlayersOrder: turnSequence,
			Deck:         player.Deck(),
		}
	})

	server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
		return game.MoveRecord{
			PlayerID:  g.CurrentPlayer().ID(),
			Timestamp: time.Now(),
			//Move: nop,
			StateDelta: g.FullState(player.ID()),
		}
	})
}

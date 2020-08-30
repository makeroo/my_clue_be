package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// MoveHandler handles vote start requests.
type MoveHandler struct{}

// RequestType returns Vote Start Request identifier.
func (*MoveHandler) RequestType() data.MessageType {
	return data.MessageVoteStartRequest
}

// BodyReader parses MoveRequest json from ws.
func (*MoveHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.MoveRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes move requests.
func (*MoveHandler) Handle(server *web.Server, req *web.Request) {
	move, ok := req.Body.(*data.MoveRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting MoveRequest, found", req.Body)
		return
	}

	g, err := server.CheckStartedGame(req.UserIO)

	if err != nil {
		req.SendError(err)

		return
	}

	record, err := g.Move(move.EnterRoom, move.MapX, move.MapY)

	if err != nil {
		req.SendError(err)

		return
	}
	/*
		if game.IsRoom(movingPlayer.Room) {
			game.History = append(game.History, MoveRecord{
				PlayerID:  movingPlayer.PlayerID,
				timestamp: time.Now(),
				move: &EnterRoomMove{
					room: movingPlayer.Room,
				},
			})
		} else {
			game.History = append(game.History, MoveRecord{
				PlayerID:  movingPlayer.PlayerID,
				timestamp: time.Now(),
				move: &MovingInTheHallwayMove{
					mapX: movingPlayer.MapX,
					mapY: movingPlayer.MapY,
				},
			})
		}

		message := NotifyGameState{
			State:          game.state,
			CurrentPlayer:  game.Players[game.currentPlayer].PlayerID,
			RemainingSteps: game.remainingSteps,

			PlayerPositions: []PlayerPosition{
				{
					PlayerID: movingPlayer.PlayerID,
					Room:     movingPlayer.Room,
					MapX:     movingPlayer.MapX,
					MapY:     movingPlayer.MapY,
				},
			},
		}

		server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
			return message
		})
	*/
	server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
		return record.AsMessageFor(player)
	})
}

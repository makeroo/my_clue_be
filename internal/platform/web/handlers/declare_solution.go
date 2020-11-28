package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// DeclareSolutionHandler handles sign in requests.
type DeclareSolutionHandler struct{}

// RequestType returns Declare Solution Request identifier.
func (*DeclareSolutionHandler) RequestType() data.MessageType {
	return data.MessageDeclareSolutionRequest
}

// BodyReader parses DeclareSolutionRequest json from ws.
func (*DeclareSolutionHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.DeclareSolutionRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes declare solution requests.
func (*DeclareSolutionHandler) Handle(server *web.Server, req *web.Request) {
	declareSolution, ok := req.Body.(*data.DeclareSolutionRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting DeclareSolutionRequest, found", req.Body)
		return
	}

	g, err := server.CheckCurrentPlayer(req)

	if err != nil {
		req.SendError(err)

		return
	}

	records, err := g.CheckSolution(declareSolution.Character, declareSolution.Room, declareSolution.Weapon)

	if err != nil {
		req.SendError(err)

		return
	}

	req.SendMessage(data.MessageEmptyResponse, nil)

	for _, record := range records {
		server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
			return record.AsMessageFor(player)
		})
	}
}

package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// QuerySolutionHandler handles vote start requests.
type QuerySolutionHandler struct{}

// RequestType returns Query Solution Request identifier.
func (*QuerySolutionHandler) RequestType() data.MessageType {
	return data.MessageQuerySolutionRequest
}

// BodyReader parses QuerySolutionRequest json from ws.
func (*QuerySolutionHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.QuerySolutionRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes query solution requests.
func (*QuerySolutionHandler) Handle(server *web.Server, req *web.Request) {
	querySolution, ok := req.Body.(*data.QuerySolutionRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting QuerySolutionRequest, found", req.Body)
		return
	}

	g, err := server.CheckCurrentPlayer(req)

	if err != nil {
		req.SendError(err)

		return
	}

	record, err := g.QuerySolution(querySolution.Character, querySolution.Weapon)

	if err != nil {
		req.SendError(err)

		return
	}

	req.SendMessage(data.MessageEmptyResponse, nil)

	server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
		return record.AsMessageFor(player)
	})
}

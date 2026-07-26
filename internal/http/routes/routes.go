package routes

import (
	"github.com/go-chi/chi/v5"

	"db-studio-go/internal/http/handlers"
)

func RegisterAPIRoutes(
	r chi.Router,
	connHandler *handlers.ConnectionHandler,
	tableHandler *handlers.TableHandler,
	queryHandler *handlers.QueryHandler,
	historyHandler *handlers.HistoryHandler,
	snippetHandler *handlers.SnippetHandler,
	ddlHandler *handlers.DDLHandler,
) {
	r.Route("/api", func(api chi.Router) {
		api.Get("/connection/status", connHandler.HandleGetConnectionStatus)
		api.Get("/tables", tableHandler.HandleGetTables)
		api.Get("/schema/graph", tableHandler.HandleGetSchemaGraph)
		api.Post("/tables", tableHandler.HandleCreateTable)
		api.Get("/tables/{name}/schema", tableHandler.HandleGetSchema)
		api.Get("/tables/{name}/data", tableHandler.HandleGetData)
		api.Post("/tables/{name}", tableHandler.HandleInsertRow)
		api.Post("/tables/{name}/batch", tableHandler.HandleBatchInsertOrUpdate)
		api.Patch("/tables/{name}", tableHandler.HandleUpdateRow)
		api.Delete("/tables/{name}", tableHandler.HandleDeleteRow)
		api.Post("/tables/{name}/columns", tableHandler.HandleAddColumn)
		api.Delete("/tables/{name}/columns/{col}", tableHandler.HandleDropColumn)
		api.Get("/tables/{name}/ddl", ddlHandler.HandleGetTableDDL)
		api.Post("/query", queryHandler.HandleExecuteQuery)

		// History endpoints
		api.Get("/history", historyHandler.HandleGetHistory)
		api.Delete("/history", historyHandler.HandleClearHistory)

		// Snippet endpoints
		api.Get("/snippets", snippetHandler.HandleGetSnippets)
		api.Post("/snippets", snippetHandler.HandleSaveSnippet)
		api.Delete("/snippets/{id}", snippetHandler.HandleDeleteSnippet)

		// DDL Export endpoint
		api.Get("/export/ddl", ddlHandler.HandleExportFullDDL)
	})
}

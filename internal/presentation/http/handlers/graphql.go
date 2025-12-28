package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/academic"
	"github.com/k/iRegistro/internal/presentation/graphql"
)

func GraphQLHandler(service *academic.AcademicService) gin.HandlerFunc {
	// Create Resolver with dependency
	resolver := &graphql.Resolver{
		AcademicService: service,
	}

	// Create Executable Schema
	executableSchema := graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver})

	// Create Server
	h := handler.NewDefaultServer(executableSchema)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

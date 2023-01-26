package graph_data_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Allocation struct {
	Name string  `json:"name"`
	Y    float64 `json:"y"`
}

func PercentageAllocations(server *server.Webserver, ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")

	allocations := server.UnitOfWork.AllocationRepository.GetByPortfolioId(portfolioId)

	var allocationsModel []Allocation
	for _, allocation := range allocations.Entries {
		allocationsModel = append(allocationsModel, Allocation{
			Name: allocation.Symbol,
			Y:    allocation.Percentage,
		})
	}
	ctx.JSON(http.StatusOK, allocationsModel)
}

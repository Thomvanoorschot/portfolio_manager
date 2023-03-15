package graph_data_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/Thomvanoorschot/portfolioManager/app/models/graph_data_models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PercentageAllocations struct {
	allocationRepository *repositories.AllocationRepository
}

func NewPercentageAllocations(allocationRepository *repositories.AllocationRepository) *PercentageAllocations {
	return &PercentageAllocations{allocationRepository: allocationRepository}
}

func (handler *PercentageAllocations) Handle(ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")

	allocations := handler.allocationRepository.GetByPortfolioId(portfolioId)

	var allocationsModel []graph_data_models.Allocation
	for _, allocation := range allocations.Entries {
		allocationsModel = append(allocationsModel, graph_data_models.Allocation{
			Name: allocation.Symbol,
			Y:    allocation.PercentageOfTotal,
		})
	}
	ctx.JSON(http.StatusOK, allocationsModel)
}

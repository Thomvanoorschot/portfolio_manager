package graph_data_handlers

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Allocation struct {
	Name string  `json:"name"`
	Y    float64 `json:"y"`
}

type PercentageAllocations struct {
	allocationRepository *repositories.AllocationRepository
}

func NewPercentageAllocations(allocationRepository *repositories.AllocationRepository) *PercentageAllocations {
	return &PercentageAllocations{allocationRepository: allocationRepository}
}

func (handler *PercentageAllocations) Handle(ctx *gin.Context) {
	portfolioId := ctx.Param("portfolioId")

	allocations := handler.allocationRepository.GetByPortfolioId(portfolioId)

	var allocationsModel []Allocation
	for _, allocation := range allocations.Entries {
		allocationsModel = append(allocationsModel, Allocation{
			Name: allocation.Symbol,
			Y:    allocation.Percentage,
		})
	}
	ctx.JSON(http.StatusOK, allocationsModel)
}

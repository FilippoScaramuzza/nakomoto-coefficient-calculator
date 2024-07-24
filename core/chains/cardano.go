package chains

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type CardanoResponse struct {
	Label string  `json:"label"`
	Epoch int     `json:"epoch"`
	Stake float64 `json:"stake"`
}

func Cardano() (int, error) {
	url := "https://www.balanceanalytics.io/api/mavdata.json"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return 0, err
	}
	defer resp.Body.Close()

	var responseData []struct {
		CardanoResponse []CardanoResponse `json:"mav_data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return 0, err
	}

	var votingPowers []big.Int
	for _, data := range responseData {
		for _, mavData := range data.CardanoResponse {
			stakeInt := big.NewInt(int64(mavData.Stake))
			votingPowers = append(votingPowers, *stakeInt)
		}
	}

	// Calculate total voting power
	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)

	// Calculate Nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums51(totalVotingPower, votingPowers)

	fmt.Println("Total voting power:", totalVotingPower)
	fmt.Println("The Nakamoto coefficient for Cardano is", nakamotoCoefficient)

	// Return Nakamoto coefficient
	return nakamotoCoefficient, nil
}
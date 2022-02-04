package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/suite"
)

type ParamsTestSuite struct {
	suite.Suite
}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}

func (suite *ParamsTestSuite) TestParamKeyTable() {
	suite.Require().IsType(paramtypes.KeyTable{}, ParamKeyTable())
}

func (suite *ParamsTestSuite) TestParamsValidate() {
	validExponentialCalculation := ExponentialCalculation{
		A: sdk.NewDec(int64(300_000_000)),
		R: sdk.NewDecWithPrec(5, 1),
		C: sdk.NewDec(int64(9_375_000)),
		B: sdk.OneDec(),
	}

	validInflationDistribution := InflationDistribution{
		StakingRewards:  sdk.NewDecWithPrec(533334, 6),
		UsageIncentives: sdk.NewDecWithPrec(333333, 6),
		CommunityPool:   sdk.NewDecWithPrec(133333, 6),
	}

	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{
			"default",
			DefaultParams(),
			false,
		},
		{
			"valid",
			NewParams(
				"aevmos",
				validExponentialCalculation,
				validInflationDistribution,
			),
			false,
		},
		{
			"valid param literal",
			Params{
				MintDenom:              "aevmos",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution:  validInflationDistribution,
			},
			false,
		},
		{
			"invalid - denom",
			NewParams(
				"/aevmos",
				validExponentialCalculation,
				validInflationDistribution,
			),
			true,
		},
		{
			"invalid - denom",
			Params{
				MintDenom:              "",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution:  validInflationDistribution,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative A",
			Params{
				MintDenom: "aevmos",
				ExponentialCalculation: ExponentialCalculation{
					A: sdk.NewDec(int64(-1)),
					R: sdk.NewDecWithPrec(5, 1),
					C: sdk.NewDec(int64(9_375_000)),
					B: sdk.OneDec(),
				},
				InflationDistribution: validInflationDistribution,
			},
			true,
		},
		{
			"invalid - exponential calculation - R greater than 1",
			Params{
				MintDenom: "aevmos",
				ExponentialCalculation: ExponentialCalculation{
					A: sdk.NewDec(int64(300_000_000)),
					R: sdk.NewDecWithPrec(5, 0),
					C: sdk.NewDec(int64(9_375_000)),
					B: sdk.OneDec(),
				},
				InflationDistribution: validInflationDistribution,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative R",
			Params{
				MintDenom: "aevmos",
				ExponentialCalculation: ExponentialCalculation{
					A: sdk.NewDec(int64(300_000_000)),
					R: sdk.NewDecWithPrec(-5, 1),
					C: sdk.NewDec(int64(9_375_000)),
					B: sdk.OneDec(),
				},
				InflationDistribution: validInflationDistribution,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative C",
			Params{
				MintDenom: "aevmos",
				ExponentialCalculation: ExponentialCalculation{
					A: sdk.NewDec(int64(300_000_000)),
					R: sdk.NewDecWithPrec(5, 1),
					C: sdk.NewDec(int64(-9_375_000)),
					B: sdk.OneDec(),
				},
				InflationDistribution: validInflationDistribution,
			},
			true,
		},
		{
			"invalid - exponential calculation - R greater than 1",
			Params{
				MintDenom: "aevmos",
				ExponentialCalculation: ExponentialCalculation{
					A: sdk.NewDec(int64(300_000_000)),
					R: sdk.NewDecWithPrec(5, 0),
					C: sdk.NewDec(int64(9_375_000)),
					B: sdk.NewDec(int64(2)),
				},
				InflationDistribution: validInflationDistribution,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative B",
			Params{
				MintDenom: "aevmos",
				ExponentialCalculation: ExponentialCalculation{
					A: sdk.NewDec(int64(300_000_000)),
					R: sdk.NewDecWithPrec(5, 1),
					C: sdk.NewDec(int64(9_375_000)),
					B: sdk.OneDec().Neg(),
				},
				InflationDistribution: validInflationDistribution,
			},
			true,
		},
		{
			"invalid - inflation distribution - negative staking rewards",
			Params{
				MintDenom:              "aevmos",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution: InflationDistribution{
					StakingRewards:  sdk.OneDec().Neg(),
					UsageIncentives: sdk.NewDecWithPrec(333333, 6),
					CommunityPool:   sdk.NewDecWithPrec(133333, 6),
				},
			},
			true,
		},
		{
			"invalid - inflation distribution - negative usage incentives",
			Params{
				MintDenom:              "aevmos",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution: InflationDistribution{
					StakingRewards:  sdk.NewDecWithPrec(533334, 6),
					UsageIncentives: sdk.OneDec().Neg(),
					CommunityPool:   sdk.NewDecWithPrec(133333, 6),
				},
			},
			true,
		},
		{
			"invalid - inflation distribution - negative community pool rewards",
			Params{
				MintDenom:              "aevmos",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution: InflationDistribution{
					StakingRewards:  sdk.NewDecWithPrec(533334, 6),
					UsageIncentives: sdk.NewDecWithPrec(333333, 6),
					CommunityPool:   sdk.OneDec().Neg(),
				},
			},
			true,
		},
		{
			"invalid - inflation distribution - total distribution ratio unequal 1",
			Params{
				MintDenom:              "aevmos",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution: InflationDistribution{
					StakingRewards:  sdk.NewDecWithPrec(533333, 6),
					UsageIncentives: sdk.NewDecWithPrec(333333, 6),
					CommunityPool:   sdk.NewDecWithPrec(133333, 6),
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			suite.Require().Error(err, tc.name)
		} else {
			suite.Require().NoError(err, tc.name)
		}
	}
}
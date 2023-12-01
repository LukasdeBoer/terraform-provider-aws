// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securityhub

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindActionTargetByARN(ctx context.Context, conn *securityhub.Client, arn string) (*types.ActionTarget, error) {
	input := &securityhub.DescribeActionTargetsInput{
		ActionTargetArns: []string{arn},
	}

	output, err := conn.DescribeActionTargets(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.ActionTargets) == 0 {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.ActionTargets); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	return &output.ActionTargets[0], nil
}

func FindFindingAggregatorByARN(ctx context.Context, conn *securityhub.Client, arn string) (*securityhub.GetFindingAggregatorOutput, error) {
	input := &securityhub.GetFindingAggregatorInput{
		FindingAggregatorArn: aws.String(arn),
	}

	output, err := conn.GetFindingAggregator(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func FindAdminAccount(ctx context.Context, conn *securityhub.Client, adminAccountID string) (*types.AdminAccount, error) {
	input := &securityhub.ListOrganizationAdminAccountsInput{}

	output, err := conn.ListOrganizationAdminAccounts(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	for _, adminAccount := range output.AdminAccounts {
		if adminAccount.AccountId == aws.String(adminAccountID) {
			return &adminAccount, nil
		}
	}

	return nil, err
}

func FindInsight(ctx context.Context, conn *securityhub.Client, arn string) (*types.Insight, error) {
	input := &securityhub.GetInsightsInput{
		InsightArns: []string{arn},
		MaxResults:  aws.Int32(1),
	}

	output, err := conn.GetInsights(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.Insights) == 0 {
		return nil, nil
	}

	return &output.Insights[0], nil
}

func FindStandardsControlByStandardsSubscriptionARNAndStandardsControlARN(ctx context.Context, conn *securityhub.Client, standardsSubscriptionARN, standardsControlARN string) (*types.StandardsControl, error) {
	input := &securityhub.DescribeStandardsControlsInput{
		StandardsSubscriptionArn: aws.String(standardsSubscriptionARN),
	}

	output, err := conn.DescribeStandardsControls(ctx, input)

	var result types.StandardsControl
	for _, control := range output.Controls {
		if aws.ToString(control.StandardsControlArn) == standardsControlARN {
			result = control
			break
		}
	}

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return &result, nil
}

package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/test"
	"testing"
)

func TestConfigs(t *testing.T) {
	configurationSchema := field.NewConfiguration(
		ConfigurationFields,
		FieldRelationships...,
	)

	testCases := []test.TestCase{
		// Add test cases here.
	}

	test.ExerciseTestCases(t, configurationSchema, ValidateConfig, testCases)
}

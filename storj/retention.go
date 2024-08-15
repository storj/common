// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

// RetentionMode represents the retention mode of an object version.
type RetentionMode uint8

const (
	// NoRetention signifies that a retention period has not been set on an object version.
	NoRetention RetentionMode = 0

	// ComplianceMode signifies that an object version is locked in compliance mode
	// and cannot be deleted or modified until the retention period expires.
	ComplianceMode RetentionMode = 1

	// GovernanceMode signifies that an object version is locked in governance mode
	// and cannot be deleted or modified until the retention period expires or the lock is removed.
	GovernanceMode RetentionMode = 3

	// LegalHold signifies that an object version is locked in legal hold
	// and cannot be deleted or modified until the legal hold is removed.
	LegalHold RetentionMode = 4

	// LegalHoldAndComplianceMode is a helper definition for combined Legal Hold and Compliance modes.
	LegalHoldAndComplianceMode = LegalHold | ComplianceMode
	// LegalHoldAndGovernanceMode is a helper definition for combined Legal Hold and Governance modes.
	LegalHoldAndGovernanceMode = LegalHold | GovernanceMode
)

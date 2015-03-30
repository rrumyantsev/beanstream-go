package fields

import ()

const (
	TransactionId     = iota + 1 // >=,<=,=,<>,>,<
	Amount                       // >=,<=,=,<>,>,<
	MaskedCardNumber             // =
	CardOwner                    // =,START WITH
	OrderNumber                  // >=,<=,=,<>,>,<
	IPAddress                    // =,START WITH
	AuthorizationCode            // =,START WITH
	TransType                    // =
	CardType                     // =
	Response                     // =
	BillingName                  // =,START WITH
	BillingEmail                 // =,START WITH
	BillingPhone                 // =,START WITH
	ProcessedBy                  // =
	Ref1                         // =,START WITH
	Ref2                         // =,START WITH
	Ref3                         // =,START WITH
	Ref4                         // =,START WITH
	Ref5                         // =,START WITH
	ProductName                  // =,START WITH
	ProductID                    // =,START WITH
	CustCode                     // =,START WITH
	IDAdjustmentTo               // =
	IDAdjustedBy                 // =
)

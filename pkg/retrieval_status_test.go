package complyancesdk

import "testing"

func TestGetDocumentStatusRequiresDocumentID(t *testing.T) {
	cfg := NewSDKConfig("test-key", EnvironmentSandbox, []*Source{}, nil)
	if err := Configure(cfg); err != nil {
		t.Fatalf("configure failed: %v", err)
	}

	_, err := GetDocumentStatus("   ")
	if err == nil {
		t.Fatalf("expected error for empty documentID")
	}
}

func TestGetSubmissionStatusIsDeprecated(t *testing.T) {
	cfg := NewSDKConfig("test-key", EnvironmentSandbox, []*Source{}, nil)
	if err := Configure(cfg); err != nil {
		t.Fatalf("configure failed: %v", err)
	}

	_, err := GetSubmissionStatus("sub-123")
	if err == nil {
		t.Fatalf("expected deprecation error for submission polling")
	}
}

func TestGetStatusAliasIsDeprecated(t *testing.T) {
	cfg := NewSDKConfig("test-key", EnvironmentSandbox, []*Source{}, nil)
	if err := Configure(cfg); err != nil {
		t.Fatalf("configure failed: %v", err)
	}

	_, err := GetStatus("sub-123")
	if err == nil {
		t.Fatalf("expected deprecation error for submission polling alias")
	}
}

func TestSubmitPayloadRequiresPayload(t *testing.T) {
	sourceType := SourceTypeFirstParty
	sources := []*Source{NewSource("src", "1", &sourceType)}
	cfg := NewSDKConfig("test-key", EnvironmentSandbox, sources, nil)
	if err := Configure(cfg); err != nil {
		t.Fatalf("configure failed: %v", err)
	}

	_, err := SubmitPayload("   ", "src:1", CountrySA, DocumentTypeTaxInvoice)
	if err == nil {
		t.Fatalf("expected validation error for empty payload")
	}
}

func TestSubmitPayloadRejectsUnknownSource(t *testing.T) {
	cfg := NewSDKConfig("test-key", EnvironmentSandbox, []*Source{}, nil)
	if err := Configure(cfg); err != nil {
		t.Fatalf("configure failed: %v", err)
	}

	_, err := SubmitPayload("{\"invoice\":\"ok\"}", "src:1", CountrySA, DocumentTypeTaxInvoice)
	if err == nil {
		t.Fatalf("expected error for unknown source")
	}
}

func TestSubmitPayloadReturnsMockedSuccessWithValidInput(t *testing.T) {
	sourceType := SourceTypeFirstParty
	sources := []*Source{NewSource("src", "1", &sourceType)}
	cfg := NewSDKConfig("test-key", EnvironmentSandbox, sources, nil)
	if err := Configure(cfg); err != nil {
		t.Fatalf("configure failed: %v", err)
	}

	response, err := SubmitPayload("{\"invoice\":\"ok\"}", "src:1", CountrySA, DocumentTypeTaxInvoice)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if response.GetSubmissionID() != "mock-id" {
		t.Fatalf("expected mock-id, got %s", response.GetSubmissionID())
	}
	if response.GetStatus() != SubmissionStatusSubmitted {
		t.Fatalf("expected SUBMITTED, got %s", response.GetStatus())
	}
}

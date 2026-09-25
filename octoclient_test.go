package octoclient

import "testing"

func TestConvertByteToStruct_ParsesNextAction(t *testing.T) {
	body := []byte(`{
		"msg": "invoked",
		"requestId": "550e8400-e29b-41d4-a716-446655440000",
		"data": {"referenceId": "REF987"},
		"nextAction": {
			"type": "invoke-service",
			"serviceCode": "fetch-account-details",
			"data": {"detailsUrl": "https://vendor.example.com/details"}
		}
	}`)

	resp, err := ConvertByteToStruct(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.NextAction == nil {
		t.Fatal("expected NextAction to be non-nil")
	}
	if resp.NextAction.Type != NextActionTypeInvokeService {
		t.Errorf("Type = %q, want %q", resp.NextAction.Type, NextActionTypeInvokeService)
	}
	if resp.NextAction.ServiceCode != "fetch-account-details" {
		t.Errorf("ServiceCode = %q, want %q", resp.NextAction.ServiceCode, "fetch-account-details")
	}
	if resp.NextAction.Data["detailsUrl"] != "https://vendor.example.com/details" {
		t.Errorf("Data[detailsUrl] = %v, want %q", resp.NextAction.Data["detailsUrl"], "https://vendor.example.com/details")
	}
}

// TestConvertByteToStruct_NoNextAction confirms older/unchanged responses
// (no nextAction key at all) still parse cleanly, with NextAction left nil.
func TestConvertByteToStruct_NoNextAction(t *testing.T) {
	body := []byte(`{"msg": "invoked", "requestId": "550e8400-e29b-41d4-a716-446655440000", "data": {}}`)

	resp, err := ConvertByteToStruct(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.NextAction != nil {
		t.Errorf("expected NextAction to be nil, got %+v", resp.NextAction)
	}
}

func TestNextAction_ToPayload(t *testing.T) {
	na := &NextAction{
		Type:        NextActionTypeInvokeService,
		ServiceCode: "fetch-account-details",
		VendorID:    "vendor-1",
		Data:        map[string]interface{}{"detailsUrl": "https://vendor.example.com/details"},
	}

	payload, err := na.ToPayload()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payload.ServiceCode != na.ServiceCode {
		t.Errorf("ServiceCode = %q, want %q", payload.ServiceCode, na.ServiceCode)
	}
	if payload.VendorID != na.VendorID {
		t.Errorf("VendorID = %q, want %q", payload.VendorID, na.VendorID)
	}
	if payload.Data["detailsUrl"] != na.Data["detailsUrl"] {
		t.Errorf("Data[detailsUrl] = %v, want %v", payload.Data["detailsUrl"], na.Data["detailsUrl"])
	}
}

// TestNextAction_ToPayload_UnrecognizedType confirms the forward-compat rule:
// an unrecognized Type must not be silently acted on.
func TestNextAction_ToPayload_UnrecognizedType(t *testing.T) {
	na := &NextAction{Type: "some-future-type"}

	_, err := na.ToPayload()
	if err == nil {
		t.Fatal("expected an error for an unrecognized nextAction type, got nil")
	}
}

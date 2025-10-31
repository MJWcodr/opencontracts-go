package opencontractinggo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	//"github.com/stretchr/testify/assert"
)

// Test data fixtures
var (
	exampleDataFolderValid = "./exampleData/valid"
	exampleDataFolderBroken = "./exampleData/broken"
	validOpenContract = OpenContract{
		URI:           "http://example.com/contract/1",
		PublishedDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Extensions:    []string{"ext1", "ext2"},
		Releases: []Release{
			{
				OCID: "ocid-123",
				ID:   "release-123",
				Date: "2023-01-01",
				Tag:  []string{"tender"},
				Tender: Tender{
					ID:          "tender-123",
					Title:       "Test Tender",
					Description: "Test Description",
					ProcuringEntity: TenderProcuringEntity{
						Name: "Test Entity",
						ID:   "entity-123",
					},
				},
				Buyer: Buyer{
					Name: "Test Buyer",
					ID:   "buyer-123",
				},
			},
		},
	}
)

func TestParseJSONToBiddingOffer(t *testing.T) {
	t.Run("Valid JSON", func(t *testing.T) {
		jsonBytes, err := json.Marshal(validOpenContract)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}

		result := ParseJSONToBiddingOffer(jsonBytes, true)
		if result.URI != validOpenContract.URI {
			t.Errorf("Expected URI %s, got %s", validOpenContract.URI, result.URI)
		}
	})
	// ...existing code for other test cases...
}

func TestParseOpenContractToJSON(t *testing.T) {
	jsonBytes := ParseOpenContractToJSON(validOpenContract)
	var parsed OpenContract
	err := json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal generated JSON: %v", err)
	}
	if parsed.URI != validOpenContract.URI {
		t.Errorf("Expected URI %s, got %s", validOpenContract.URI, parsed.URI)
	}
}

func TestSynthesizeTenderID(t *testing.T) {
	release := Release{
		ID:     "release-123",
		Tender: Tender{
			// Empty ID field
		},
	}
	t.Run("Successfully synthesize ID", func(t *testing.T) {
		tender := SynthesizeTenderID(release)
		expected := release.ID + "-tender"
		if tender.ID != expected {
			t.Errorf("Expected tender ID %s, got %s", expected, tender.ID)
		}
	})
}

func TestParseJSONFileToBiddingOffer(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "transform_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	jsonBytes, err := json.Marshal(validOpenContract)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	filePath := filepath.Join(tmpDir, "test_contract.json")
	err = os.WriteFile(filePath, jsonBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	result := ParseJSONFileToBiddingOffer(filePath)
	if result.URI != validOpenContract.URI {
		t.Errorf("Expected URI %s, got %s", validOpenContract.URI, result.URI)
	}
}

func TestAddressStruct(t *testing.T) {
	addr := Address{
		Street:      "123 Main St",
		Locality:    "Anytown",
		Region:      "State",
		PostalCode:  "12345",
		Country:     "Country",
		Description: "Description",
	}
	jsonBytes, err := json.Marshal(addr)
	if err != nil {
		t.Fatalf("Failed to marshal Address struct: %v", err)
	}
	var parsed Address
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal Address struct: %v", err)
	}
	if !reflect.DeepEqual(addr, parsed) {
		t.Errorf("Address serialization/deserialization failed. Expected %+v, got %+v", addr, parsed)
	}
}

func TestComplexStructSerialization(t *testing.T) {
	tender := Tender{
		ID:          "tender-123",
		Title:       "Complex Test Tender",
		Description: "A complex tender with nested structures",
		ProcuringEntity: TenderProcuringEntity{
			Name: "Procuring Entity",
			ID:   "entity-123",
		},
		Items: []TenderItem{
			{
				ID: "item-1",
				Classification: TenderItemClassification{
					Scheme:      "CPV",
					ID:          "12345678",
					Description: "Item Classification",
				},
				Lots: []TenderItemLot{
					{
						ID:          "lot-1",
						Title:       "Lot 1",
						Description: "First lot",
					},
				},
				RelatedLot: "lot-1",
				DeliveryAddress: Address{
					Street:     "Delivery Street",
					Locality:   "Delivery Town",
					PostalCode: "54321",
				},
			},
		},
	}
	jsonBytes, err := json.Marshal(tender)
	if err != nil {
		t.Fatalf("Failed to marshal complex Tender struct: %v", err)
	}
	var parsed Tender
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal complex Tender struct: %v", err)
	}
	if parsed.ID != tender.ID || parsed.Title != tender.Title {
		t.Errorf("Complex structure top-level fields mismatch. Expected ID=%s, Title=%s, got ID=%s, Title=%s",
			tender.ID, tender.Title, parsed.ID, parsed.Title)
	}
}

func TestExampleData(t *testing.T){
	// List all files
	files, err := os.ReadDir(exampleDataFolderValid)
	if err != nil {
		t.Fatal(err)
	}

	// Decode Files
	for _, file := range files {
		fileName := exampleDataFolderValid + "/" + file.Name() 
		f, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatal(err)
		}
		ParseJSONToBiddingOffer(f, true)
	}	
}


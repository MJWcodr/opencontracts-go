package opencontractinggo

import (
	"bytes"
	"encoding/json"
	"os"
	"slices"

	"github.com/charmbracelet/log"
)

type Address struct {
	Street      string `json:"streetAddress"`
	Locality    string `json:"locality"`
	Region      string `json:"region"`
	PostalCode  string `json:"postalCode"`
	Country     string `json:"countryName"`
	Description string `json:"description"`
}

type ContactPoint struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	FaxNumber string `json:"faxNumber"`
	Url       string `json:"url"`
}

type Identifier struct {
	Id        string `json:"id"`
	LegalName string `json:"legalName"`
	Scheme    string `json:"scheme"`
}

type Party struct {
	Name         string       `json:"name"`
	ID           string       `json:"id"`
	Identifier   Identifier   `json:"identifier"`
	Address      Address      `json:"address"`
	ContactPoint ContactPoint `json:"contactPoint"`
	Roles        []string     `json:"roles"`
}

type Buyer struct {
	Name         string       `json:"name"`
	ID           string       `json:"id"`
	Identifier   Identifier   `json:"identifier"`
	Address      Address      `json:"address"`
	ContactPoint ContactPoint `json:"contactPoint"`
}

type Tender struct {
	ID                         string                `json:"id"`
	Title                      string                `json:"title"`
	Description                string                `json:"description"`
	ProcuringEntity            TenderProcuringEntity `json:"procuringEntity"`
	Items                      []TenderItem          `json:"items"`
	Awards                     []TenderAward         `json:"awards"`
	Language                   string                `json:"language"`
	ProcurementMethodDetails   string                `json:"procurementMethodDetails"`
	MainProcurementCategory    string                `json:"mainProcurementCategory"`
	NumberOfTenderers          int                   `json:"numberOfTenderers"`
	Documents                  []Document            `json:"documents"`
	ReleaseOcid                string
	Status                     string        `json:"status"`
	MinValue                   MonetaryValue `json:"minValue"`
	Value                      MonetaryValue `json:"value"`
	ProcurementMethod          string        `json:"procurementMethod"`
	ProcurementMethodRationale string        `json:"procurementMethodRationale"`
	AwardCriteria              string        `json:"awardCriteria"`
	AwardCriteriaDetails       string        `json:"awardCriteriaDetails"`
	SubmissionMethod           []string      `json:"submissionMethod"`
	SubmissionMethodDetails    string        `json:"submissionMethodDetails"`
	TenderPeriod               Period        `json:"tenderPeriod"`
	EnquiryPeriod              Period        `json:"enquiryPeriod"`
	AwardPeriod                Period        `json:"awardPeriod"`
	ContractPeriod             Period        `json:"contractPeriod"`
	HasEnquiries               bool          `json:"hasEnquiries"`
	Lots                       []any         `json:"lots" extensions:lots` // Lots extension
}

type MonetaryValue struct {
	Amount   float64   `json:"amount"`
	Currency string `json:"currency"`
}

type Period struct {
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	DurationInDays int    `json:"durationInDays"`
}

type Document struct {
	ID            string   `json:"id"`
	URL           string   `json:"url"`
	RelatedLots   []string `json:"relatedLots"`
	DocumentType  string   `json:"documentType"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	DatePublished string   `json:"datePublished"`
	Format        string   `json:"format"`
	Language      string   `json:"language"`
}

type TenderProcuringEntity struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type TenderAward struct {
	ID          string   `json:"id"`
	Status      string   `json:"status"`
	RelatedLots []string `json:"relatedLots"`
}

type DeliverLocation struct {
	Description string `json:"description" extension:"deliveryLocation"`
}

type TenderItem struct {
	ID                        string                   `json:"id"`
	Classification            TenderItemClassification `json:"classification"`
	AdditionalClassifications []TenderItemClassification
	Lots                      []TenderItemLot `json:"lots"`
	RelatedLot                string          `json:"relatedLot"`
	DeliveryAddress           Address         `json:"deliveryAddress"`
	Description               string          `json:"description"`
	Quantity                  int             `json:"quantity"`
	Unit                      TenderUnit      `json:"unit"`
	DeliveryLocation 					DeliverLocation `json:"deliveryLocation"`	
}

type TenderItemClassification struct {
	Scheme      string `json:"scheme"`
	ID          string `json:"id"`
	Description string `json:"description"`
	URI         string `json:"uri"`
}

type TenderUnit struct {
	Name   string        `json:"name"`
	Id     string        `json:"id"`
	Scheme string        `json:"scheme"`
	Value  MonetaryValue `json:"value"`
}

type TenderItemLot struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Release struct {
	OCID              string    `json:"ocid"`
	ID                string    `json:"id"`
	Date              string    `json:"date"`
	Tag               []string  `json:"tag"`
	Parties           []Party   `json:"parties"`
	Buyer             Buyer     `json:"buyer"`
	Tender            Tender    `json:"tender"`
	Publisher         Publisher `json:"publisher"`
	License           string    `json:"license"`
	PublicationPolicy string    `json:"publicationPolicy"`
	Language          string    `json:"language"`
	InitiationType    string    `json:"initiationType"`
}

type Publisher struct {
	Name   string `json:"name"`
	Scheme string `json:"scheme"`
	Uid    string `json:"uid"`
	Uri    string `json:"uri"`
}

type OpenContract struct {
	Version           string    `json:"version"`
	URI               string    `json:"uri"`
	PublishedDate     string    `json:"publishedDate"`
	Extensions        []string  `json:"extensions"`
	Releases          []Release `json:"releases"`
	Publisher         Publisher `json:"publisher"`
	License           string    `json:"license"`
	PublicationPolicy string    `json:"publicationPolicy"`
}

func SynthesizeTenderID(release Release) Tender {
	tender := release.Tender

	// check if the tender ID is already set
	if tender.ID != "" {
		panic("This should only be called when the tender ID is not set")
	}

	// Check if the release ID is set
	if release.ID == "" {
		panic("Releases.ID is required")
	}

	// set the tender ID
	tender.ID = release.ID + "-tender"

	return tender
}

func SynthesizeTenderName(release Release) Tender {
	tender := release.Tender

	// Check if the tender name is already set
	if tender.Title != "" {
		panic("This should only be called when the tender name is not set")
	}

	// Check if the release ID is set
	if release.ID == "" {
		panic("Releases.ID is required")
	}

	// set the tender name
	tender.Title = release.ID + "-tender"

	return tender
}

func (oc *OpenContract) checkOCDSVersion() error {
	// To be implemented
	return nil
}

// checkExtensions checks for what extensions are used.
// Based on the extensions additional fields may be availible in the OpenContract struct.
// Since Go does not support dynamic fields, this function will just return an error if a field is defined
// which is added by an extension that is not listed in the Extensions field.
func (oc *OpenContract) checkExtensions() error {
	supportedExtensions := map[string]string{
		"Lots": "https://raw.githubusercontent.com/open-contracting-extensions/ocds_lots_extension/bd6a5a0617abcdcf9da19db0e70a899e6b10c91c/extension.json",
	}

	if !slices.Contains(oc.Extensions, supportedExtensions["Lots"]) {
		// Check if the fields are not empty
		// If they are not empty, return an error
		// Since the Extension is not loaded.

	}


	return nil
}

func ParseJSONToBiddingOffer(byteArray []byte, strict bool) OpenContract {
	var openContract OpenContract

	
	reader := bytes.NewBuffer(byteArray)
	decoder := json.NewDecoder(reader)

	if strict {
		decoder.DisallowUnknownFields()
	}

	err := decoder.Decode(&openContract)
	if err != nil {
		panic(err)
	}

	// check for required fields
	if openContract.URI == "" {
		panic("URI is required")
	}
	if openContract.Releases[0].ID == "" {
		panic("Releases.ID is required")
	}
	if openContract.Releases[0].Tender.ID == "" {
		openContract.Releases[0].Tender.ID = SynthesizeTenderID(openContract.Releases[0]).ID
		log.Debug("Tender.ID was not set, synthesized it")
	}
	if openContract.Releases[0].Tender.Title == "" {
		log.Debug("Tender.Title was not set, synthesized it")
		openContract.Releases[0].Tender.Title = SynthesizeTenderName(openContract.Releases[0]).Title
	}
	if openContract.Releases[0].Tender.ProcuringEntity.ID == "" {
		log.Debug("Currently not using ProcuringEntity.ID")
	}
	if openContract.Releases[0].Tender.ProcuringEntity.Name == "" {
		log.Debug("Currently not using ProcuringEntity.Name")
	}
	if openContract.Releases[0].Buyer.ID == "" {
		log.Debug("Currently not using Buyer.ID")
	}

	return openContract
}

func ParseJSONFileToBiddingOffer(filePath string) OpenContract {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return ParseJSONToBiddingOffer(bytes, false)
}

func ParseOpenContractToJSON(openContract OpenContract) []byte {
	bytes, err := json.Marshal(openContract)
	if err != nil {
		panic(err)
	}
	return bytes
}

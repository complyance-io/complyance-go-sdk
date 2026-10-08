// Package complyance identifies the Complyance Go SDK module.
//
// Import github.com/complyance-io/complyance-go-sdk/v3/pkg for the APIClient,
// configured SDK functions, request models, and raw XML/JSON submission methods.
// APIClient.SubmitXmlSingle accepts UTF-8 XML bytes; SubmitXmlSingleFile reads
// an XML file unchanged. SubmitJSONSingle requires an explicit contract selector;
// SubmitLegacyJSONSingle omits it. The raw methods require a public API origin
// and send once without queueing or retrying.
package complyance

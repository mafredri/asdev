// +build !go1.9

// Code generated by cdpgen. DO NOT EDIT.

package dom

import (
	"github.com/mafredri/cdp/protocol"
)

// Node DOM interaction is implemented in terms of mirror objects that represent the actual DOM nodes. DOMNode is a base node mirror type.
type Node struct {
	NodeID NodeID `json:"nodeId"` // Node identifier that is passed into the rest of the DOM messages as the nodeId. Backend will only push node with given id once. It is aware of all requested nodes and will only fire DOM events for nodes known to the client.
	// ParentID The id of the parent node if any.
	//
	// Note: This property is experimental.
	ParentID *NodeID `json:"parentId,omitempty"`
	// BackendNodeID The BackendNodeId for this node.
	//
	// Note: This property is experimental.
	BackendNodeID  BackendNodeID `json:"backendNodeId"`
	NodeType       int           `json:"nodeType"`                 // Node's nodeType.
	NodeName       string        `json:"nodeName"`                 // Node's nodeName.
	LocalName      string        `json:"localName"`                // Node's localName.
	NodeValue      string        `json:"nodeValue"`                // Node's nodeValue.
	ChildNodeCount *int          `json:"childNodeCount,omitempty"` // Child count for Container nodes.
	Children       []Node        `json:"children,omitempty"`       // Child nodes of this node when requested with children.
	Attributes     []string      `json:"attributes,omitempty"`     // Attributes of the Element node in the form of flat array [name1, value1, name2, value2].
	DocumentURL    *string       `json:"documentURL,omitempty"`    // Document URL that Document or FrameOwner node points to.
	// BaseURL Base URL that Document or FrameOwner node uses for URL completion.
	//
	// Note: This property is experimental.
	BaseURL        *string        `json:"baseURL,omitempty"`
	PublicID       *string        `json:"publicId,omitempty"`       // DocumentType's publicId.
	SystemID       *string        `json:"systemId,omitempty"`       // DocumentType's systemId.
	InternalSubset *string        `json:"internalSubset,omitempty"` // DocumentType's internalSubset.
	XMLVersion     *string        `json:"xmlVersion,omitempty"`     // Document's XML version in case of XML documents.
	Name           *string        `json:"name,omitempty"`           // Attr's name.
	Value          *string        `json:"value,omitempty"`          // Attr's value.
	PseudoType     PseudoType     `json:"pseudoType,omitempty"`     // Pseudo element type for this node.
	ShadowRootType ShadowRootType `json:"shadowRootType,omitempty"` // Shadow root type.
	// FrameID Frame ID for frame owner elements.
	//
	// Note: This property is experimental.
	FrameID         *protocol.PageFrameID `json:"frameId,omitempty"`
	ContentDocument *Node                 `json:"contentDocument,omitempty"` // Content document for frame owner elements.
	// ShadowRoots Shadow root list for given element host.
	//
	// Note: This property is experimental.
	ShadowRoots []Node `json:"shadowRoots,omitempty"`
	// TemplateContent Content document fragment for template elements.
	//
	// Note: This property is experimental.
	TemplateContent *Node `json:"templateContent,omitempty"`
	// PseudoElements Pseudo elements associated with this node.
	//
	// Note: This property is experimental.
	PseudoElements   []Node `json:"pseudoElements,omitempty"`
	ImportedDocument *Node  `json:"importedDocument,omitempty"` // Import document for the HTMLImport links.
	// DistributedNodes Distributed nodes for given insertion point.
	//
	// Note: This property is experimental.
	DistributedNodes []BackendNode `json:"distributedNodes,omitempty"`
	// IsSVG Whether the node is SVG.
	//
	// Note: This property is experimental.
	IsSVG *bool `json:"isSVG,omitempty"`
}

// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sandpiper

/* this model is not currently used (but keeping for web socket sync */

// SyncRequest models the sync communication request
type SyncRequest struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Slice   Slice  `json:"slice"`
}

// SyncResponse models the sync communication request
type SyncResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Slice   Slice  `json:"slice"`
}

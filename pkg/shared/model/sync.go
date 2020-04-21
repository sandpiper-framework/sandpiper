// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

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

/* Actions:
slices: return slices I'm subscribed to and are available for sync

*/

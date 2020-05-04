// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper_test

import (
	"testing"

	"sandpiper/pkg/shared/model"
)

func TestChangePassword(t *testing.T) {
	user := &sandpiper.User{
		FirstName: "TestGuy",
	}

	hashedPassword := "h4$h3D"

	user.ChangePassword(hashedPassword)
	if user.PasswordChanged.IsZero() {
		t.Errorf("Last password change was not changed")
	}

	if user.Password != hashedPassword {
		t.Errorf("Password was not changed")
	}
}

func TestUpdateLastLogin(t *testing.T) {
	user := &sandpiper.User{
		FirstName: "TestGuy",
	}

	token := "helloWorld"

	user.UpdateLastLogin(token)
	if user.LastLogin.IsZero() {
		t.Errorf("Last login time was not changed")
	}

	if user.Token != token {
		t.Errorf("Token was not changed")
	}
}

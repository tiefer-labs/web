// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build race

package server

// raceEnabled is true under -race, which slows code down several times,
// so timing budgets are only checked without it.
const raceEnabled = true

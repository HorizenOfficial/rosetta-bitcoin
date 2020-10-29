// Copyright (c) 2013-2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package zenutil

const (
	//ZentoshiPerZencent is the number of zentoshi in one zencoin cent.
	ZentoshiPerZencent = 1e6

	// ZentoshiPerZen is the number of zentoshi in one zen (1 ZEN).
	ZentoshiPerZen = 1e8

	// MaxZentoshi is the maximum transaction amount allowed in zentoshi.
	MaxZentoshi = 21e6 * ZentoshiPerZen
)

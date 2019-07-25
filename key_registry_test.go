/*
 * Copyright 2019 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package badger

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildRegistry(t *testing.T) {
	storageKey := make([]byte, 32)
	dir, err := ioutil.TempDir("", "badger-test")
	_, err = rand.Read(storageKey)
	require.NoError(t, err)
	kr, err := OpenKeyRegistry(dir, false, storageKey)
	defer os.Remove(dir)
	require.NoError(t, err)
	dk, err := kr.getDataKey()
	require.NoError(t, err)
	kr.lastCreated = 0
	dk1, err := kr.getDataKey()
	require.NoError(t, err)
	require.NoError(t, kr.Close())
	kr2, err := OpenKeyRegistry(dir, false, storageKey)
	require.NoError(t, err)
	require.Equal(t, 2, len(kr2.dataKeys))
	require.Equal(t, dk.Data, kr.dataKeys[dk.KeyId].Data)
	require.Equal(t, dk1.Data, kr.dataKeys[dk1.KeyId].Data)
	require.NoError(t, kr2.Close())
}

func TestRewriteRegistry(t *testing.T) {
	dir, err := ioutil.TempDir("", "badger-test")
	require.NoError(t, err)
	storageKey := make([]byte, 32)
	_, err = rand.Read(storageKey)
	require.NoError(t, err)
	kr, err := OpenKeyRegistry(dir, false, storageKey)
	defer os.Remove(dir)
	require.NoError(t, err)
	_, err = kr.getDataKey()
	require.NoError(t, err)
	kr.lastCreated = 0
	_, err = kr.getDataKey()
	require.NoError(t, err)
	require.NoError(t, kr.Close())
	delete(kr.dataKeys, 1)
	require.NoError(t, RewriteRegistry(dir, kr, storageKey))
	kr2, err := OpenKeyRegistry(dir, false, storageKey)
	require.NoError(t, err)
	require.Equal(t, 1, len(kr2.dataKeys))
	require.NoError(t, kr2.Close())
}
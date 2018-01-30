/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package consistenthash

import (
	"fmt"
	"testing"
)

var(
	hashTestData = HashNew()
)

func TestHashAdd(t *testing.T) {
	for i := 0; i < 5; i++ {
		host := fmt.Sprintf("host%v", i+1)
		hashTestData.Add(10, host)
	}
	for k, v := range hashTestData.KeyMap {
		t.Logf("k:%v v:%v", k, v)
	}
	//t.Logf("SrcKeys:%v", hashTestData.SrcKeys)
	t.Logf("Keys:%v", hashTestData.KeyAlive)
}

func TestHashDel(t *testing.T) {
	for i := 0; i < 5; i++ {
		host := fmt.Sprintf("host%v", i+1)
		hashTestData.Add(10, host)
	}
	for k, v := range hashTestData.KeyMap {
		t.Logf("k:%v v:%v", k, v)
	}
	//t.Logf("SrcKeys:%v", hashTestData.SrcKeys)
	t.Logf("Keys:%v", hashTestData.KeyAlive)

	hashTestData.Update("host1", false)
	t.Logf("---------")
	for k, v := range hashTestData.KeyMap {
		t.Logf("k:%v v:%v", k, v)
	}
	//t.Logf("SrcKeys:%v", hashTestData.SrcKeys)
	t.Logf("Keys:%v", hashTestData.KeyAlive)
}

func TestMD5(t *testing.T) {
	for i:=0; i<10; i++{
		key := MD5(fmt.Sprintf("hello%v", i))
		t.Logf(key)
	}
}

func TestHashGet(t *testing.T) {
	for i := 0; i < 5; i++ {
		host := fmt.Sprintf("host%v", i+1)
		hashTestData.Add(10, host)
	}
	rangeLimit := 1000
	dieHost := "host2"
	prefix := "/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM%05d"
	t.Logf("---------1 正常情况--------")
	pieMap := map[string] int {
		"host1": 0,
		"host2": 0,
		"host3": 0,
		"host4": 0,
		"host5": 0,
	}
	for i:=0; i<rangeLimit; i++ {
		key := fmt.Sprintf(prefix, i)
		host, err := hashTestData.Get(key)
		if err != nil {
			t.Fatal(err)
		}
		pieMap[host] += 1
	}
	t.Logf("Keys:%v", getMapKeys(hashTestData.KeyAlive))
	t.Logf("pieMap: %v", pieMap)

	hashTestData.Update(dieHost, false)
	t.Logf("-----2. 转移: %v-----数量:%v----", dieHost, pieMap[dieHost])
	t.Logf("Keys:%v", getMapKeys(hashTestData.KeyAlive))
	pieMap2 := map[string] int {
		"host1": 0,
		"host2": 0,
		"host3": 0,
		"host4": 0,
		"host5": 0,
	}
	for i:=0; i<rangeLimit; i++ {
		key := fmt.Sprintf(prefix, i)
		host, err := hashTestData.Get(key)
		if err != nil {
			t.Fatal(err)
		}
		pieMap2[host] += 1
	}
	diffMap := map[string] int {
		"host1": 0,
		"host2": 0,
		"host3": 0,
		"host4": 0,
		"host5": 0,
	}
	for k, _ := range pieMap {
		diffMap[k] = pieMap2[k] - pieMap[k]
	}
	t.Logf("pieMap1: %v", pieMap)
	t.Logf("pieMap2: %v", pieMap2)
	t.Logf("diffMap: %v", diffMap)

	hashTestData.Update(dieHost, true)
	t.Logf("-----3. 恢复: %v-----数量:%v----", dieHost, pieMap[dieHost])
	t.Logf("Keys:%v", getMapKeys(hashTestData.KeyAlive))
	pieMap3 := map[string] int {
		"host1": 0,
		"host2": 0,
		"host3": 0,
		"host4": 0,
		"host5": 0,
	}
	for i:=0; i<rangeLimit; i++ {
		key := fmt.Sprintf(prefix, i)
		host, err := hashTestData.Get(key)
		if err != nil {
			t.Fatal(err)
		}
		pieMap3[host] += 1
	}
	t.Logf("pieMap3: %v", pieMap3)
}
/*
host2: 由正常->丢失->正常

rangelimit = 1000 时，测试效果:
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host3 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00001 new_key:63679c4b91ea4317dc91c590c4e4a9f6 idx:1 idxMove:1
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host5 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00009 new_key:0d53b9c3d8f936cbb34388f564102bcd idx:1 idxMove:3
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host3 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00019 new_key:e23b6778382e61f7f4a2c6d13376c23f idx:1 idxMove:1
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host4 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00020 new_key:5f05925eddaa552cff9bf4c1021cf587 idx:1 idxMove:2
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host3 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00021 new_key:ade6da88e87f8371abfbb7942df342a1 idx:1 idxMove:1
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host4 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00023 new_key:89adeeec2a6390a8f41ce8ee87b075e5 idx:1 idxMove:2
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host4 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00024 new_key:bc17bf4ce0c6a001b365fcd14b134771 idx:1 idxMove:2
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host3 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00028 new_key:348e91305211292eaf79da3d48a81a29 idx:1 idxMove:1
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host3 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00035 new_key:e672022fd1d2748f3ec41e521f1b0580 idx:1 idxMove:1
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host4 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00036 new_key:9961393674b8a4e94e84255bc6bec09e idx:1 idxMove:2
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host4 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00038 new_key:785e9a720fe7311ed3c6dd687391cf87 idx:1 idxMove:2
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host3 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00042 new_key:933498daecfc505c0f7044360300570d idx:1 idxMove:1
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host3 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00057 new_key:bbafb53d577f9752bfc52134507b8c8a idx:1 idxMove:1
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host5 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00058 new_key:55a92b442a07d04f6638cbc6788856e0 idx:1 idxMove:3
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host4 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00059 new_key:0dd23db51b423d3c03f55af41b2b5559 idx:1 idxMove:2
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host1 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00064 new_key:48c613c7281a03d828fee9edbfb95b76 idx:1 idxMove:0
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host4 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00074 new_key:5f4831d37c531d78166484203abb2895 idx:1 idxMove:2
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host4 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00075 new_key:37f959c819ce14cf31fd88e22a69a560 idx:1 idxMove:2
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host1 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00081 new_key:e39e5bf8ac2fee77f97d37e3074c415e idx:1 idxMove:0
2018/01/26 15:37:14 [D] [consistenthash.go:106] srcHost:host2->moveHost:host1 key:/b/apk/a29tLm1vYmlsZS5sZWdlbmRzXzExNTIxMzMxX2U4ZGIzOTM00084 new_key:ce0d8fefb6223d18a698dbb4d63e3cd0 idx:1 idxMove:0
	consistenthash_test.go:58 : ---------1 正常情况--------
	consistenthash_test.go:74 : Keys:[host1 host2 host3 host4 host5]
	consistenthash_test.go:75 : pieMap: map[host1:23 host2:20 host3:18 host4:21 host5:18]
	consistenthash_test.go:78 : -----2. 转移: host2-----数量:20----
	consistenthash_test.go:79 : Keys:[host1 host3 host4 host5]
	consistenthash_test.go:105: pieMap1: map[host4:21 host5:18 host1:23 host2:20 host3:18]
	consistenthash_test.go:106: pieMap2: map[host1:26 host2:0 host3:25 host4:29 host5:20]
	consistenthash_test.go:107: diffMap: map[host3:7 host4:8 host5:2 host1:3 host2:-20]
	consistenthash_test.go:110: -----3. 恢复: host2-----数量:20----
	consistenthash_test.go:111: Keys:[host1 host2 host3 host4 host5]
	consistenthash_test.go:127: pieMap3: map[host1:23 host2:20 host3:18 host4:21 host5:18]

rangelimit = 1000 时，测试效果:
consistenthash_test.go:58: ---------1 正常情况--------
	consistenthash_test.go:74: Keys:[host1 host2 host3 host4 host5]
	consistenthash_test.go:75: pieMap: map[host1:202 host2:212 host3:203 host4:171 host5:212]
	consistenthash_test.go:78: -----2. 转移: host2-----数量:212----
	consistenthash_test.go:79: Keys:[host1 host3 host4 host5]
	consistenthash_test.go:105: pieMap1: map[host3:203 host4:171 host5:212 host1:202 host2:212]
	consistenthash_test.go:106: pieMap2: map[host3:258 host4:236 host5:257 host1:249 host2:0]
	consistenthash_test.go:107: diffMap: map[host1:47 host2:-212 host3:55 host4:65 host5:45]
	consistenthash_test.go:110: -----3. 恢复: host2-----数量:212----
	consistenthash_test.go:111: Keys:[host1 host2 host3 host4 host5]
	consistenthash_test.go:127: pieMap3: map[host1:202 host2:212 host3:203 host4:171 host5:212]
*/

